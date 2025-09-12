// Package wasm provides a thin Go wrapper around a wasm-bindgen compiled WebAssembly module
// used by this repository. It offers helpers to instantiate the module with wazero,
// manage guest memory, pass strings safely, and decode the common (ptr,len) and
// return-area patterns used by the Rust/WASM boundary.
package wasm

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"

	"biscuit-wasm-go/utils"
)

const defaultWasmRelPath = "target/wasm32-unknown-unknown/release/biscuit_wasm_go.wasm"

// WasmEnv wraps a wazero Module and its context, exposing convenience methods to
// find exports and call them
type WasmEnv struct {
	Ctx        context.Context
	Module     api.Module
	returnArea uint64
}

// Call calls the specified exported function with the specified parameters.
func (env WasmEnv) Call(name string, params ...uint64) ([]uint64, error) {
	function := env.Module.ExportedFunction(name)
	if function == nil {
		slog.Error("exported function not found", slog.String("name", name))
		return nil, fmt.Errorf("exported function '%s' not found", name)
	}
	data, err := function.Call(env.Ctx, params...)
	if err != nil {
		slog.Error("failed to call exported function", slog.String("name", name), slog.Any("err", err))
		return nil, err
	}
	return data, nil
}

// GetReturnArea returns the address of 12 bytes allocated memory area.
// This area is used to deserialize the return value from the WASM function.
// The area is allocated once and reused for all calls.
// Warning: This method doesn't ensure that the area hasn't been freed.
func (env WasmEnv) GetReturnArea() (uint64, error) {

	// If already allocated, return it
	if env.returnArea != 0 {
		return env.returnArea, nil
	}

	// Allocate memory
	returnArea, err := env.Call("get_return_area")
	if err != nil {
		return 0, err
	}

	return returnArea[0], err
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

// InitWasm loads, compiles, and instantiates the Biscuit WASM module using wazero.
// It also installs host import stubs so the module can run without a JS host.
func InitWasm() (WasmEnv, error) {
	ctx := context.Background()
	// Create a new runtime
	runtimeConfig := wazero.NewRuntimeConfig().WithMemoryCapacityFromMax(true).WithDebugInfoEnabled(true)
	runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)

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

	wasmEnv := WasmEnv{
		Ctx:        ctx,
		Module:     module,
		returnArea: 0,
	}

	return wasmEnv, nil

}

// malloc allocates memory in the WASM memory.
func (env WasmEnv) malloc(ln uint64) (uint64, error) {

	dataPtr, err := env.Call("malloc", ln, 1)
	if err != nil {
		return 0, err
	}

	return dataPtr[0], nil
}

// Free frees memory in the WASM memory.
func (env WasmEnv) Free(ptr uint64, ln uint64) {

	_, err := env.Call("free", ptr, ln, 1)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to free memory at %x", ptr), slog.Any("err", err))

	}

}

// WriteBytesToWasm write writes data to the WASM memory.
func (env WasmEnv) WriteBytesToWasm(data []byte) (uint64, error) {

	ptr, err := env.malloc(uint64(len(data)))
	if err != nil {
		panic(err)
	}

	mem := env.Module.Memory()
	ok := mem.Write(uint32(ptr), data)
	if !ok {
		panic("failed to write data")
	}
	return ptr, nil
}
