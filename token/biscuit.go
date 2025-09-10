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
	return "biscuit_to_base64"
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

	returnArea, err := env.GetReturnArea()
	if err != nil {
		return Biscuit{}, err
	}

	// Write the base64 string into wasm memory and pass (ptr, len)
	strPtr, err := env.WriteBytesToWasm([]byte(biscuitBase64))
	if err != nil {
		return Biscuit{}, err
	}
	defer env.Free(strPtr, uint64(len(biscuitBase64)))

	_, err = env.Call("biscuit_from_base64", returnArea, strPtr, uint64(len(biscuitBase64)), publicKey.Ptr())
	if err != nil {
		return Biscuit{}, fmt.Errorf("biscuit_fromBase64 failed: %w", err)
	}

	ptr, err := env.ResultPointer(returnArea)

	if err != nil {
		fmt.Println(err)
		return Biscuit{}, err
	}

	biscuit.env = env
	biscuit.ptr = ptr
	return biscuit, nil
}

// ToBytes returns the byte slice representation of the Biscuit.
func (biscuit Biscuit) ToBytes() ([]byte, error) {
	if biscuit.ptr == 0 {
		return nil, fmt.Errorf("token not initialized")
	}

	returnArea, err := biscuit.env.GetReturnArea()
	if err != nil {
		return nil, err
	}

	_, err = biscuit.env.Call("biscuit_to_bytes", returnArea, biscuit.ptr)
	if err != nil {
		return nil, err
	}

	return biscuit.env.ResultBytes(returnArea)

}

// FromBytes creates a new Biscuit from a byte slice.
func (biscuit Biscuit) FromBytes(env wasm.WasmEnv, bytes []byte, publicKey keypair.PublicKey) (Biscuit, error) {

	returnArea, err := env.GetReturnArea()
	if err != nil {
		return Biscuit{}, err
	}

	// Write the base64 string into wasm memory and pass (ptr, len)
	strPtr, err := env.WriteBytesToWasm(bytes)
	if err != nil {
		return Biscuit{}, err
	}
	defer env.Free(strPtr, uint64(len(bytes)))

	_, err = env.Call("biscuit_from_bytes", returnArea, strPtr, uint64(len(bytes)), publicKey.Ptr())
	if err != nil {
		return Biscuit{}, fmt.Errorf("biscuit_fromBase64 failed: %w", err)
	}

	ptr, err := env.ResultPointer(returnArea)

	if err != nil {
		return Biscuit{}, err
	}

	biscuit.env = env
	biscuit.ptr = ptr
	return biscuit, nil

}
