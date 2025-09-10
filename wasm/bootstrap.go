// Package wasm contains host-side bootstrap for the wazero runtime. It auto-instantiates
// minimal host functions to provide entropy for wasm module.
package wasm

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// InstantiateImportStubs builds host modules to satisfy wasm-bindgen imports.
func InstantiateImportStubs(ctx context.Context, runtime wazero.Runtime, compiledModule wazero.CompiledModule) error {
	imports := compiledModule.ImportedFunctions()

	builders := map[string]wazero.HostModuleBuilder{}
	for _, def := range imports {
		modName, name, isImport := def.Import()
		if !isImport {
			continue
		}

		builder, ok := builders[modName]
		if !ok {
			builder = runtime.NewHostModuleBuilder(modName)
			builders[modName] = builder
		}
		params := def.ParamTypes()
		results := def.ResultTypes()

		switch name {
		case "print":
			builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, module api.Module, stack []uint64) {
				mem := module.Memory()
				ln := api.DecodeU32(stack[1])
				buf, ok := mem.Read(uint32(stack[0]), ln)
				if !ok {
					panic("failed to read data")
				}
				// print to stdout debug wasm message
				fmt.Println(string(buf))
				stack[0] = 0

			}), params, results).Export(name)

		case "__getrandom_custom":
			builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, module api.Module, stack []uint64) {

				mem := module.Memory()

				arr := api.DecodeU32(stack[0])
				ln := uint32(stack[1])

				buf := make([]byte, ln)
				if n, err := rand.Read(buf); err == nil {
					if uint32(n) < ln {
						for i := n; uint32(i) < ln; i++ {
							buf[i] = 0
						}
					}
					_ = mem.Write(arr, buf)
				}

				stack[0] = 0

			}), params, results).Export(name)

		default:
			builder.NewFunctionBuilder().WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
				// no-op
			}), params, results).Export(name)

		}
	}

	for modName, b := range builders {
		if _, err := b.Instantiate(ctx); err != nil {
			panic(fmt.Sprintf("failed to instantiate host module %q: %w", modName, err))
		}
	}
	return nil
}
