// Package shared contains cross-package helpers for interacting with the WASM boundary.
package shared

import (
	"biscuit-wasm-go/wasm"
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

	_, err = env.Call(function, resultPtr, value.Ptr())
	if err != nil {
		slog.Error("biscuitbuilder_toString failed", slog.Any("err", err))
		return "", err
	}

	return env.GetStringValueFromPointer(resultPtr)

}
