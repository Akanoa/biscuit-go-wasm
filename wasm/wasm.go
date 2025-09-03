// Package wasm provides a thin Go wrapper around a wasm-bindgen compiled WebAssembly module
// used by this repository. It offers helpers to instantiate the module with wazero,
// manage guest memory, pass strings safely, and decode the common (ptr,len) and
// return-area patterns used by the Rust/WASM boundary.
package wasm

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"

	"biscuit-wasm-go/utils"
)

const defaultWasmRelPath = "target/wasm32-unknown-unknown/release/biscuit_wasm_go.wasm"

// WasmEnv wraps a wazero Module and its context, exposing convenience methods to
// find exports and call them, and to manage the common memory patterns used by
// wasm-bindgen generated interfaces.
type WasmEnv struct {
	Ctx    context.Context
	Module api.Module
}

func (env WasmEnv) GetFunction(name string) (api.Function, error) {
	function := env.Module.ExportedFunction(name)
	if function == nil {
		slog.Error("exported function not found", slog.String("name", name))
		return nil, fmt.Errorf("exported function '%s' not found", name)
	}
	return function, nil
}

func (env WasmEnv) GetMemory() (api.Memory, error) {
	memory := env.Module.Memory()
	if memory == nil {
		return nil, fmt.Errorf("exported memory '%s' not found", "default")
	}
	return memory, nil
}

func (env WasmEnv) Call(function api.Function, params ...uint64) ([]uint64, error) {
	return function.Call(env.Ctx, params...)
}

func CloseRuntime(runtime wazero.Runtime, ctx context.Context) {
	if runtime.Close(ctx) != nil {

		panic("failed to close runtime")
	}
}

func CloseWasmModule(module api.Module, goContext context.Context) {
	if module.Close(goContext) != nil {
		panic("failed to close module")
	}
}

// InitWasm loads, compiles and instantiates the Biscuit WASM module using wazero.
// It also installs host import stubs so the module can run without a JS host.
func InitWasm() (WasmEnv, error) {
	ctx := context.Background()
	// Create a new runtime
	runtime := wazero.NewRuntime(ctx)

	var sourceWasm []byte
	var err error

	// Resolve path dynamically
	path, err := utils.ResolveAssetFile(defaultWasmRelPath)
	if err != nil {
		slog.Error("Unable to resolve wasm file", slog.Any("err", err))
		return WasmEnv{}, err
	}

	sourceWasm, err = os.ReadFile(path)
	if err != nil {
		slog.Error("Unable to read wasm file", slog.String("file", path), slog.Any("err", err))
		return WasmEnv{}, err
	}

	// Compile module
	compiled, err := runtime.CompileModule(ctx, sourceWasm)
	if err != nil {
		slog.Error("Unable to compile wasm file", slog.String("file", path), slog.Any("err", err))
		return WasmEnv{}, err
	}

	// Auto-instantiate host stubs for any imported functions (e.g., from "__wbindgen_placeholder__").
	if err := InstantiateImportStubs(ctx, runtime, compiled); err != nil {
		slog.Error("Unable to instantiate import stubs", slog.Any("err", err))
		return WasmEnv{}, err
	}

	// Use default module config so the module's start function (if any) runs.
	wasmConfig := wazero.NewModuleConfig()

	module, err := runtime.InstantiateModule(ctx, compiled, wasmConfig)
	if err != nil {
		slog.Error("Unable to instantiate module", slog.Any("err", err))
		return WasmEnv{}, err
	}

	return WasmEnv{
		Ctx:    ctx,
		Module: module,
	}, nil
}

func (env WasmEnv) Free(ptr uint64, length uint64) error {
	free, err := env.GetFunction("__wbindgen_free")
	if err != nil {
		slog.Error("exported function not found", slog.String("name", "__wbindgen_free"))
		return err
	}
	_, err = env.Call(free, ptr, length, 1)
	return err
}

func (env WasmEnv) Malloc(length uint64) (uint64, error) {
	malloc, err := env.GetFunction("__wbindgen_malloc")
	if err != nil {
		slog.Error("exported function not found", slog.String("name", "__wbindgen_malloc"))
		return 0, err
	}
	results, err := env.Call(malloc, length, 1)
	if err != nil {
		slog.Error("malloc failed", slog.Any("err", err))
		return 0, err
	}

	if len(results) != 1 {
		slog.Error("malloc failed: unexpected return value")
		return 0, fmt.Errorf("malloc failed: unexpected return value")
	}

	return results[0], nil
}

// GetStringValueFromPointer string is a double-pointed value. The first pointer is a pointer to the return area,
// ptr pointed to an 8-byte area with the following layout:
// 0: 4 bytes: string pointer
// 4: 4 bytes: string length
// This second pointer is the actual string data, we read the length and decode the string from memory
// and free the return area.
//
// Memory Layout Diagram:
// +----------------+     +-------------------+
// | Return Area    |     | String Data      |
// | (8 bytes)      |     | (variable length)|
// +----------------+     +-------------------+
// | String Ptr   --|---->| Actual string    |
// | String Length  |     | content...       |
// +----------------+     +-------------------+
//
//	^
//	|
//
// ptr (input parameter)
func (env WasmEnv) GetStringValueFromPointer(ptr uint64) (string, error) {

	// read return area
	mem := env.Module.Memory()
	buf, ok := mem.Read(uint32(ptr), 8)
	if !ok {
		slog.Error("cannot read return area")
		return "", fmt.Errorf("cannot read return area")
	}
	strPtr := binary.LittleEndian.Uint32(buf[0:4])
	strLen := binary.LittleEndian.Uint32(buf[4:8])

	// decode string from memory
	strBytes, ok := mem.Read(strPtr, strLen)
	if !ok {
		panic("cannot read string")
	}
	stringData := string(strBytes)

	err := env.Free(uint64(strPtr), uint64(strLen))
	if err != nil {
		slog.Error("cannot free string", slog.Uint64("ptr", uint64(strPtr)), slog.Uint64("len", uint64(strLen)))
		return "", err
	}

	return stringData, nil
}

// GetError retrieves the error associated with a given externref index from the ExternrefTableMirror.
// It returns a string representation of the error or an empty string if the error cannot be resolved.
// An error is returned if the provided index is invalid.
func (env WasmEnv) GetError(idx uint64) (string, error) {
	if int(idx) >= len(ExternrefTableMirror) {
		return "", fmt.Errorf("unknown error: invalid externref index %d", idx)
	}

	v := ExternrefTableMirror[idx]
	switch data := v.(type) {
	case nil:
		return "unknown error", nil
	case string:
		return data, nil
	case map[string]any:
		// Prefer a stable, human-friendly serialization without Go's map[...] prefix
		// Common shape from wasm-bindgen is nested single-key maps representing error enums, e.g.:
		// {"FailedLogic": {"NoMatchingPolicy": {"checks": {}}}}
		// Collapse nested single-key maps into a path like "FailedLogic: NoMatchingPolicy".
		var parts []string
		cur := data
		for {
			if len(cur) != 1 {
				break
			}
			var k string
			var v any
			for kk, vv := range cur { k, v = kk, vv }
			parts = append(parts, k)
			// descend if the value is another map[string]any
			next, ok := v.(map[string]any)
			if !ok {
				// If value is an empty map[any]any or prints as map[], stop and emit path
				if fmt.Sprintf("%v", v) == "map[]" {
					break
				}
				// Non-map leaf: include its printable form as final segment and stop
				if v != nil {
					parts = append(parts, fmt.Sprintf("%v", v))
				}
				break
			}
			cur = next
		}
		if len(parts) > 0 {
			// Special-case to avoid redundant trailing technical segments like "checks" when empty
			if len(parts) >= 2 && parts[len(parts)-1] == "checks" {
				parts = parts[:len(parts)-1]
			}
			// Join with ": " for readability
			msg := parts[0]
			for i := 1; i < len(parts); i++ {
				msg += ": " + parts[i]
			}
			return msg, nil
		}
		// Fallback: stable order by keys
		keys := make([]string, 0, len(data))
		for k := range data {
			keys = append(keys, k)
		}
		// simple insertion sort to avoid importing sort
		for i := 1; i < len(keys); i++ {
			for j := i; j > 0 && keys[j-1] > keys[j]; j-- {
				keys[j-1], keys[j] = keys[j], keys[j-1]
			}
		}
		out := ""
		for i, k := range keys {
			if i > 0 { out += ", " }
			out += fmt.Sprintf("%s: %v", k, data[k])
		}

		return out, nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// WriteString writes a UTF-8 encoded string to the WebAssembly module's memory and returns a pointer to its location.
func (env WasmEnv) WriteString(data string) (uint64, error) {

	mem := env.Module.Memory()

	// Prepare UTF-8 bytes from data
	bytes := []byte(data)
	// Allocate buffer for string bytes
	strPtr, err := env.Malloc(uint64(len(bytes)))
	if err != nil {
		return 0, fmt.Errorf("malloc for string failed: %w", err)
	}

	// Write bytes into memory
	if ok := mem.Write(uint32(strPtr), bytes); !ok {
		_ = env.Free(strPtr, uint64(len(bytes)))
		return 0, fmt.Errorf("cannot write string bytes to wasm memory")
	}

	return strPtr, nil
}

// getArea allocates a new area of the given size in bytes.
func (env WasmEnv) getArea(size uint64) (uint64, error) {
	retPtr, err := env.Malloc(size)
	if err != nil {
		return 0, fmt.Errorf("malloc for area failed: %w", err)
	}

	return retPtr, nil
}

// ReturnAreaSize is the size of the return area used by functions that return
// Result<T, E> where wasm-bindgen lays out three u32 slots:
// 0:4 bytes: value pointer (or 0)
// 4:4 bytes: error externref index (or 0)
// 8:4 bytes: is_err (0 = Ok, non-zero = Err)
// Note: Some exports use a compact, 2-slot layout: (ptr_or_err, is_err).
// For those, use SmallReturnAreaSize and GetResult2.
const ReturnAreaSize = uint64(16)

// GetReturnArea ReturnAreaSize is the size of the return area in bytes.
func (env WasmEnv) GetReturnArea() (uint64, error) {
	// Allocate return area (3 u32 values: value_ptr, error_ptr, is_err)
	return env.getArea(ReturnAreaSize)
}

// SmallReturnAreaSize is used by some wasm-bindgen exports that encode
// Result-like values in 2 u32 slots:
// 0:4 bytes: value pointer or error externref index
// 4:4 bytes: is_err (0 = Ok, non-zero = Err)
const SmallReturnAreaSize = uint64(8)

// GetSmallReturnArea allocates a 2-slot return area (ptr_or_err, is_err).
func (env WasmEnv) GetSmallReturnArea() (uint64, error) {
	return env.getArea(SmallReturnAreaSize)
}

// StringAreaSize is the size of the string area in bytes.
// 0:4 bytes: string pointer
// 4:8 bytes: string length
const StringAreaSize = uint64(8)

// GetStringArea StringAreaSize is the size of the string area in bytes.
func (env WasmEnv) GetStringArea() (uint64, error) {
	// Allocate return area (2 u32 values: string_ptr, string_len)
	return env.getArea(StringAreaSize)
}

// GetPointee returns the value pointed to by the given pointer.
// The pointer is expected to point to a return area.
func (env WasmEnv) GetPointee(ptr uint64) (uint64, error) {
	mem := env.Module.Memory()

	// Read result triple
	buf, ok := mem.Read(uint32(ptr), uint32(ReturnAreaSize))
	if !ok {
		return 0, fmt.Errorf("cannot read return area")
	}
	valuePtr := binary.LittleEndian.Uint32(buf[0:4])
	errPtr := binary.LittleEndian.Uint32(buf[4:8])
	isErr := int32(binary.LittleEndian.Uint32(buf[8:12]))

	if isErr != 0 {
		serr, err := env.GetError(uint64(errPtr))
		if err != nil {
			return 0, fmt.Errorf("cannot get error string: %w", err)
		}
		return 0, errors.New(serr)
	}

	return uint64(valuePtr), nil
}

// GetResult2 decodes a 2-slot Result area allocated with SmallReturnAreaSize.
// It returns the value pointer on Ok, or an error constructed from the externref
// index on Err.
func (env WasmEnv) GetResult2(ptr uint64) (uint64, error) {
	mem := env.Module.Memory()
	buf, ok := mem.Read(uint32(ptr), uint32(SmallReturnAreaSize))
	if !ok {
		return 0, fmt.Errorf("cannot read small return area")
	}
	ptrOrErr := binary.LittleEndian.Uint32(buf[0:4])
	isErr := binary.LittleEndian.Uint32(buf[4:8])
	if isErr != 0 {
		msg, err := env.GetError(uint64(ptrOrErr))
		if err != nil {
			return 0, fmt.Errorf("cannot get error string: %w", err)
		}
		return 0, errors.New(msg)
	}
	return uint64(ptrOrErr), nil
}
