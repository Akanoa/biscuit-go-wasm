// Package shared contains cross-package helpers for interacting with the WASM boundary.
package shared

import (
	"biscuit-wasm-go/wasm"
	"fmt"
	"log/slog"
)

// AsString calls a wasm export corresponding to value.ToStringWasmFunction() and
// decodes the (ptr,len) result as a Go string. The value must provide a valid Ptr().
func AsString(env wasm.WasmEnv, value Stringable) (string, error) {

	wasmFunction := value.ToStringWasmFunction()

	function, err := env.GetFunction(wasmFunction)
	if err != nil {
		return "", err
	}

	resultPtr, err := env.GetStringArea()
	if err != nil {
		return "", err
	}
	defer env.Free(resultPtr, wasm.StringAreaSize)

	ret, err := env.Call(function, value.Ptr())
	if err != nil {
		slog.Error("biscuitbuilder_toString failed", slog.Any("err", err))
		return "", err
	}

	mem, _ := env.GetMemory()

	buf, ok := mem.Read(uint32(ret[0]), uint32(ret[1]))
	if !ok {
		return "", fmt.Errorf("cannot read return area")
	}

	return string(buf), nil

}
