// Package wasm contains host-side bootstrap for the wazero runtime. It auto-instantiates
// minimal host functions to satisfy imports produced by wasm-bindgen glue, so the
// compiled module can run without a JavaScript environment.
package wasm

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

const OffsetJsidx = 128

// taLen maps a synthesized typed-array handle (we use the byte offset as the handle)
// to its length. This lets entropy functions and copy helpers know where and how
// many bytes to read/write in guest memory.
var taLen = map[uint32]uint32{}

// ExternrefTableSize tracks the logical size of the wasm-bindgen externref table when hosted in Go.
var ExternrefTableSize uint32

// ExternrefTableMirror mirrors the wasm-bindgen externref table so Go code can inspect entries.
// Index 0 is reserved (undefined), and init seeds [undefined, null, true, false] similar to the JS glue.
var ExternrefTableMirror []any

var mapString = map[string]uint64{}

var memoryObjHandle uint32

var globalObjHandle uint32

type JsNull struct{}

// instantiateImportStubs builds host modules to satisfy wasm-bindgen imports.
func InstantiateImportStubs(ctx context.Context, runtime wazero.Runtime, c wazero.CompiledModule) error {
	imports := c.ImportedFunctions()
	if len(imports) == 0 {
		return nil
	}

	builders := map[string]wazero.HostModuleBuilder{}
	for _, def := range imports {
		modName, name, isImport := def.Import()
		if !isImport {
			continue
		}
		//if modName != "__wbindgen_placeholder__" && modName != "__wbindgen_externref_xform__" && modName != "wbg" {
		//	return fmt.Errorf("unsupported import module: %s.%s", modName, name)
		//}
		builder, ok := builders[modName]
		if !ok {
			builder = runtime.NewHostModuleBuilder(modName)
			builders[modName] = builder
		}
		params := def.ParamTypes()
		results := def.ResultTypes()

		switch name {
		case "__wbindgen_init_externref_table":
			builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
				fmt.Println("######################init externref table")
				if len(ExternrefTableMirror) == 0 {
					ExternrefTableMirror = append(ExternrefTableMirror, nil)
				}
				offset := uint32(len(ExternrefTableMirror))
				for i := 0; i < OffsetJsidx+4; i++ {
					ExternrefTableMirror = append(ExternrefTableMirror, nil)
				}
				ExternrefTableMirror[offset+0+OffsetJsidx] = nil
				ExternrefTableMirror[offset+1+OffsetJsidx] = JsNull{}
				ExternrefTableMirror[offset+2+OffsetJsidx] = true
				ExternrefTableMirror[offset+3+OffsetJsidx] = false
				ExternrefTableSize = uint32(len(ExternrefTableMirror))
				_ = stack
			}), params, results).Export(name)
		case "__wbg_randomFillSync_ac0988aba3254290", "__wbg_getRandomValues_b8f5dbd5f3995a9e":
			fn := api.GoModuleFunc(func(ctx context.Context, m api.Module, stack []uint64) {
				mem := m.Memory()
				_ = api.DecodeU32(stack[0])
				arr := api.DecodeU32(stack[1])
				ln := taLen[arr]
				if ln == 0 {
					return
				}
				buf := make([]byte, ln)
				if n, err := rand.Read(buf); err == nil {
					if uint32(n) < ln {
						for i := n; uint32(i) < ln; i++ {
							buf[i] = 0
						}
					}
					_ = mem.Write(arr, buf)
				}
			})
			builder.NewFunctionBuilder().WithGoModuleFunction(fn, params, results).Export(name)
		case "__wbindgen_copy_to_typed_array":
			fn := api.GoModuleFunc(func(ctx context.Context, m api.Module, stack []uint64) {
				mem := m.Memory()
				_ = api.DecodeU32(stack[0])
				srcLen := api.DecodeU32(stack[1])
				dstPtr := api.DecodeU32(stack[2])
				if srcLen == 0 {
					return
				}
				buf := make([]byte, srcLen)
				if n, err := rand.Read(buf); err == nil {
					if uint32(n) < srcLen {
						for i := n; uint32(i) < srcLen; i++ {
							buf[i] = 0
						}
					}
					_ = mem.Write(dstPtr, buf)
				}
			})
			builder.NewFunctionBuilder().WithGoModuleFunction(fn, params, results).Export(name)
		case "__wbindgen_is_object", "__wbindgen_is_function":
			builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
				stack[0] = api.EncodeU32(1)
			}), params, results).Export(name)
		case "__wbindgen_string_new":
			// Convert Rust string (ptr,len) to a JS handle (index) by storing it in the externref table mirror.
			builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, m api.Module, stack []uint64) {
				mem := m.Memory()
				ptr := api.DecodeU32(stack[0])
				ln := api.DecodeU32(stack[1])
				var s string
				idxString := uint64(len(ExternrefTableMirror)) + OffsetJsidx

				// If the string is empty, return 0.
				if ln == 0 {
					stack[0] = api.EncodeU32(0)
					return
				}

				// Read string from memory.
				buf, ok := mem.Read(ptr, ln)

				// If the read failed, return 0.
				if !ok {
					fmt.Println("failed to read string from memory")
					return
				}

				// calculate hash of string
				h := sha256.Sum256(buf)
				key := fmt.Sprintf("---- %x", h)

				// Ensure table has reserved 0 entry.
				if len(ExternrefTableMirror) == 0 {
					ExternrefTableMirror = append(ExternrefTableMirror, nil)
				}

				// search for string in map
				if idx, ok := mapString[key]; ok {
					if int(idx) < len(ExternrefTableMirror) && idx >= OffsetJsidx {
						s = ExternrefTableMirror[idx].(string)
						idxString = idx
					} else {
						// Index corruption detected, create new entry
						s = string(buf)
						idxString = uint64(len(ExternrefTableMirror)) + OffsetJsidx
						ExternrefTableMirror = append(ExternrefTableMirror, s)
						mapString[key] = idxString
					}
				} else {
					s = string(buf)
					idxString = uint64(len(ExternrefTableMirror)) + OffsetJsidx
					ExternrefTableMirror = append(ExternrefTableMirror, s)
					mapString[key] = idxString

				}

				ExternrefTableSize = uint32(uint64(len(ExternrefTableMirror)))

				stack[0] = idxString
			}), params, results).Export(name)
		//case "__wbg_require_60cc747a6bc5215a":
		//	builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
		//		stack[0] = api.EncodeExternref(1)
		//	}), params, results).Export(name)
		//case "__wbg_newwithbyteoffsetandlength_d97e637ebe145a9a":
		//	builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
		//		byteOffset := api.DecodeU32(stack[1])
		//		length := api.DecodeU32(stack[2])
		//		taLen[byteOffset] = length
		//		stack[0] = api.EncodeU32(byteOffset)
		//	}), params, results).Export(name)
		case "__wbg_newwithbyteoffsetandlength_d97e637ebe145a9a":
			builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
				byteOffset := api.DecodeU32(stack[1])
				length := api.DecodeU32(stack[2])
				taLen[byteOffset] = length
				stack[0] = api.EncodeU32(byteOffset)
			}), params, results).Export(name)
		case "__wbg_set_65595bdd868b3009":
			builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, m api.Module, stack []uint64) {
				mem := m.Memory()
				srcHandle := api.DecodeU32(stack[1])
				dstPtr := api.DecodeU32(stack[2])
				ln := taLen[srcHandle]
				if ln == 0 {
					return
				}
				// Check memory bounds before operations
				memSize := mem.Size()
				if srcHandle >= memSize || dstPtr >= memSize {
					fmt.Printf("WARNING: __wbg_set memory bounds check failed - srcHandle:%d dstPtr:%d memSize:%d\n", srcHandle, dstPtr, memSize)
					return
				}
				if srcHandle+ln > memSize || dstPtr+ln > memSize {
					fmt.Printf("WARNING: __wbg_set buffer overflow check failed - operation would exceed memory bounds\n")
					return
				}
				if buf, ok := mem.Read(srcHandle, ln); ok {
					if !mem.Write(dstPtr, buf) {
						fmt.Printf("WARNING: __wbg_set memory write failed\n")
					}
				} else {
					fmt.Printf("WARNING: __wbg_set memory read failed\n")
				}
			}), params, results).Export(name)
		case "__wbg_subarray_aa9065fa9dc5df96":
			builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
				base := api.DecodeU32(stack[0])
				begin := api.DecodeU32(stack[1])
				end := api.DecodeU32(stack[2])
				newHandle := base + begin
				var l uint32
				if end >= begin {
					l = end - begin
				}
				taLen[newHandle] = l
				stack[0] = api.EncodeU32(newHandle)
			}), params, results).Export(name)
		case "__wbg_stack_0ed75d68575b0f3c":
			// Extract stack trace from Error object and write to WASM memory
			builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, m api.Module, stack []uint64) {
				fmt.Println("****************wasm-bindgen stack")

				mem := m.Memory()
				retPtr := api.DecodeU32(stack[0])
				errorHandle := api.DecodeU32(stack[1])

				fmt.Println("wasm-bindgen stack:", errorHandle)

				// Get stack trace from error object in externref table
				stackTrace := "Error stack trace not available"
				adjustedHandle := errorHandle - OffsetJsidx
				if adjustedHandle < uint32(len(ExternrefTableMirror)) {
					if errObj := ExternrefTableMirror[adjustedHandle]; errObj != nil {
						if errMap, ok := errObj.(map[string]any); ok {
							if stackVal, exists := errMap["stack"]; exists {
								if stackStr, ok := stackVal.(string); ok {
									stackTrace = stackStr
								}
							}
						}
						// If no stack property, create a basic stack trace
						if stackTrace == "Error stack trace not available" {
							stackTrace = "Error: " + fmt.Sprintf("%v", errObj)
						}
					}
				}

				// Allocate memory for the string and write it
				strBytes := []byte(stackTrace)
				strLen := uint32(len(strBytes))

				// Find available memory space (simple allocation)
				var strPtr uint32 = 0x10000 // Start at a safe offset
				memSize := mem.Size()
				for strPtr+strLen >= memSize {
					strPtr += 0x1000 // Move to next page boundary
					if strPtr > memSize {
						strPtr = memSize - strLen
						break
					}
				}

				// Write string to memory
				if strPtr+strLen <= memSize {
					mem.Write(strPtr, strBytes)
				}

				// Write pointer and length back to the result location
				if retPtr+8 <= memSize {
					// Write pointer at retPtr+0, length at retPtr+4
					mem.WriteUint32Le(retPtr, strPtr)
					mem.WriteUint32Le(retPtr+4, strLen)
				}
				fmt.Println("wasm-bindgen stack:", stackTrace)

			}), params, results).Export(name)
		case "__wbg_new_8a6f238a6ece86ea":
			// Create a new Error object and return its handle
			builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
				// Create a new error object in the externref table
				if len(ExternrefTableMirror) == 0 {
					ExternrefTableMirror = append(ExternrefTableMirror, nil)
				}
				// Create a new error object with stack property
				errorObj := map[string]any{
					"name":    "Error",
					"message": "",
					"stack":   "Error\n    at <anonymous>",
				}
				ExternrefTableMirror = append(ExternrefTableMirror, errorObj)
				handle := uint32(len(ExternrefTableMirror)-1) + OffsetJsidx
				stack[0] = api.EncodeU32(handle)
			}), params, results).Export(name)
		case "__wbindgen_memory":
			builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
				if memoryObjHandle == 0 {
					if len(ExternrefTableMirror) == 0 {
						ExternrefTableMirror = append(ExternrefTableMirror, nil)
					}
					ExternrefTableMirror = append(ExternrefTableMirror, map[string]any{"__kind": "memory"})
					memoryObjHandle = uint32(len(ExternrefTableMirror)-1) + OffsetJsidx
				}
				stack[0] = api.EncodeU32(memoryObjHandle)
			}), params, results).Export(name)
		case "__wbindgen_throw":
			// wasm-bindgen uses this to throw JS exceptions. Do NOT panic; capture the message in the externref table
			// so callers can retrieve it via GetError without trapping the runtime.
			builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, m api.Module, stack []uint64) {
				fmt.Println("****************wasm-bindgen throw")

				mem := m.Memory()
				ptr := api.DecodeU32(stack[0])
				ln := api.DecodeU32(stack[1])
				msg := "wasm-bindgen throw"
				if ln > 0 {
					if buf, ok := mem.Read(ptr, ln); ok {
						msg = string(buf)
					}
				}

				fmt.Println("wasm-bindgen throw:", msg)

				// Store the message as a new externref entry for later retrieval by host code.
				if len(ExternrefTableMirror) == 0 {
					ExternrefTableMirror = append(ExternrefTableMirror, nil)
				}
				ExternrefTableMirror = append(ExternrefTableMirror, msg)
				// Do not panic: simply return to let the guest continue or handle error paths.
			}), params, results).Export(name)
		case "__wbg_static_accessor_SELF_37c5d418e4bf5819":
			// Static accessor for 'self' global object
			builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
				if globalObjHandle == 0 {
					if len(ExternrefTableMirror) == 0 {
						ExternrefTableMirror = append(ExternrefTableMirror, nil)
					}
					ExternrefTableMirror = append(ExternrefTableMirror, map[string]any{"__kind": "global"})
					globalObjHandle = uint32(len(ExternrefTableMirror)-1) + OffsetJsidx
				}
				stack[0] = api.EncodeU32(globalObjHandle)
			}), params, results).Export(name)
		case "__wbindgen_object_clone_ref":
			// Clone an object reference - return the same handle since Go manages memory automatically
			builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
				handle := api.DecodeU32(stack[0])
				// In a JavaScript environment, this would increment a reference count
				// In Go, we just return the same handle since Go's GC handles object lifetime
				stack[0] = api.EncodeU32(handle)
			}), params, results).Export(name)
		default:
			fmt.Println("unsupported import:", modName, name)
			if modName == "__wbindgen_externref_xform__" {
				switch name {
				case "__wbindgen_externref_table_grow":
					builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
						n := api.DecodeU32(stack[0])
						if n > 10000 { // Prevent excessive allocation
							fmt.Printf("WARNING: externref_table_grow requested %d entries, limiting to 10000\n", n)
							n = 10000
						}
						prev := uint32(len(ExternrefTableMirror))
						for i := uint32(0); i < n; i++ {
							ExternrefTableMirror = append(ExternrefTableMirror, nil)
						}
						ExternrefTableSize = uint32(len(ExternrefTableMirror))
						stack[0] = api.EncodeU32(prev + OffsetJsidx)
					}), params, results).Export(name)
				case "__wbindgen_externref_table_set_null":
					builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
						idx := api.DecodeU32(stack[0])
						actualIdx := idx + OffsetJsidx
						if actualIdx < uint32(len(ExternrefTableMirror)) && actualIdx >= OffsetJsidx {
							ExternrefTableMirror[actualIdx] = JsNull{}
						} else {
							fmt.Printf("WARNING: externref_table_set_null invalid index %d (actual %d, table size %d)\n", idx, actualIdx, len(ExternrefTableMirror))
						}
					}), params, results).Export(name)
				default:
					builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
						fmt.Println("externref_xform", name)
						_ = stack
					}), params, results).Export(name)
				}
				continue
			}
			// Default no-op for other imports (object_drop_ref, clone_ref, etc.)
			builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
				fmt.Println("default import", name)
				_ = stack
			}), params, results).Export(name)
		}
	}

	// Some externref call helpers expected on module "wbg".
	if _, ok := builders["wbg"]; !ok {
		builders["wbg"] = runtime.NewHostModuleBuilder("wbg")
	}
	b := builders["wbg"]
	b.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) { _ = stack }), []api.ValueType{api.ValueTypeExternref, api.ValueTypeExternref, api.ValueTypeExternref}, []api.ValueType{api.ValueTypeExternref}).Export("__wbg_call_7cccdd69e0791ae2")
	b.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) { _ = stack }), []api.ValueType{api.ValueTypeExternref, api.ValueTypeExternref}, []api.ValueType{api.ValueTypeExternref}).Export("__wbg_call_672a4d21634d4a24")

	for modName, b := range builders {
		if _, err := b.Instantiate(ctx); err != nil {
			return fmt.Errorf("failed to instantiate host module %q: %w", modName, err)
		}
	}
	return nil
}
