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

	returnArea, err := env.GetReturnArea()
	if err != nil {
		return "", err
	}

	_, err = env.Call(wasmFunction, returnArea, value.Ptr())
	if err != nil {
		slog.Error("biscuitbuilder_toString failed", slog.Any("err", err))
		return "", err
	}

	return env.ResultString(returnArea)

}
