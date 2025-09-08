// Package keypair provides types and helpers to work with public/private keys
// for Biscuit, bridging Go calls to the underlying WebAssembly functions.
package keypair

import (
	"biscuit-wasm-go/shared"
	"biscuit-wasm-go/wasm"
	"fmt"
)

// PublicKey represents a WASM public key.
type PublicKey struct {
	env wasm.WasmEnv
	ptr uint64
}

// ToStringWasmFunction returns the name of the WASM function that converts the PublicKey to its string representation.
func (publicKey PublicKey) ToStringWasmFunction() string {
	return "public_key_to_hex"
}

// Ptr returns the pointer to the underlying WASM object.
func (public_key PublicKey) Ptr() uint64 {
	return public_key.ptr
}

// FromString initializes a PublicKey from a string and a given SignatureAlgorithm using a Wasm environment.
// Returns the initialized PublicKey or an error in case of failure.
func (publicKey PublicKey) FromString(env wasm.WasmEnv, data string, algorithm SignatureAlgorithm) (PublicKey, error) {

	returnArea, err := env.GetReturnArea()
	if err != nil {
		return PublicKey{}, err
	}

	strPtr, err := env.WriteBytesToWasm([]byte(data))
	if err != nil {
		return PublicKey{}, err
	}

	// Call: publickey_fromString(out_ptr, str_ptr, str_len)
	_, err = env.Call("public_key_from_hex", returnArea, strPtr, uint64(len(data)), uint64(algorithm))

	ptr, err := env.ResultPointer(returnArea)
	if err != nil {
		return PublicKey{}, err
	}

	publicKey.ptr = ptr
	publicKey.env = env
	return publicKey, nil
}

// ToString converts the PublicKey to its string representation using the linked Wasm environment. Returns the string or an error.
func (publick_key PublicKey) ToString() (string, error) {
	if publick_key.ptr == 0 {
		return "", fmt.Errorf("token publick_key not initialized")
	}

	return shared.AsString(publick_key.env, publick_key)
}
