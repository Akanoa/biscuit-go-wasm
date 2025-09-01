// Package token exposes high-level Biscuit token types backed by the underlying WASM module.
// It provides convenience helpers to create tokens from base64 and to serialize them back.
package token

import (
	"biscuit-wasm-go/keypair"
	"biscuit-wasm-go/shared"
	"biscuit-wasm-go/wasm"
	"fmt"
)

// Biscuit represents a Biscuit object.
type Biscuit struct {
	env wasm.WasmEnv
	ptr uint64
}

// New creates a new Biscuit using a Wasm environment and a pointer to the WASM object.
func (biscuit Biscuit) New(env wasm.WasmEnv, ptr uint64) Biscuit {
	return Biscuit{env: env, ptr: ptr}
}

// Ptr returns the pointer to the WASM object.
func (biscuit Biscuit) Ptr() uint64 {
	return biscuit.ptr
}

// ToStringWasmFunction returns the name of the WASM function used to convert the Biscuit to a string.
func (biscuit Biscuit) ToStringWasmFunction() string {
	return "biscuit_toBase64"
}

// ToBase64 returns the string representation of the Biscuit.
func (biscuit Biscuit) ToBase64() (string, error) {
	if biscuit.ptr == 0 {
		return "", fmt.Errorf("token not initialized")
	}

	return shared.AsString(biscuit.env, biscuit)
}

// FromBase64 creates a new Biscuit from a base64 string.
func (biscuit Biscuit) FromBase64(env wasm.WasmEnv, biscuitBase64 string, publicKey keypair.PublicKey) (Biscuit, error) {

	function, err := env.GetFunction("biscuit_fromBase64")
	if err != nil {
		return Biscuit{}, err
	}

	// Write the base64 string into wasm memory and pass (ptr, len)
	strPtr, err := env.WriteString(biscuitBase64)
	if err != nil {
		return Biscuit{}, fmt.Errorf("cannot write base64 to wasm memory: %w", err)
	}
	defer env.Free(strPtr, uint64(len(biscuitBase64)))

	returnPtr, err := env.GetReturnArea()
	if err != nil {
		return Biscuit{}, err
	}
	defer env.Free(returnPtr, wasm.ReturnAreaSize)

	_, err = env.Call(function, returnPtr, strPtr, uint64(len(biscuitBase64)), publicKey.Ptr())
	if err != nil {
		return Biscuit{}, fmt.Errorf("biscuit_fromBase64 failed: %w", err)
	}

	valuePtr, err := env.GetPointee(returnPtr)

	if err != nil {
		return Biscuit{}, err
	}

	biscuit.env = env
	biscuit.ptr = valuePtr
	return biscuit, nil
}
