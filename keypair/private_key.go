package keypair

import (
	"biscuit-wasm-go/shared"
	"biscuit-wasm-go/wasm"
	"fmt"
)

// PrivateKey represents a WASM private key.
type PrivateKey struct {
	env wasm.WasmEnv
	ptr uint64
}

// ToStringWasmFunction returns the name of the WASM function that converts the PublicKey to its string representation.
func (privateKey PrivateKey) ToStringWasmFunction() string {
	return "private_key_to_hex"
}

// Ptr returns the pointer to the underlying WASM object.
func (privateKey PrivateKey) Ptr() uint64 {
	return privateKey.ptr
}

// New creates a new PrivateKey using the linked Wasm environment.
func (privateKey PrivateKey) New(env wasm.WasmEnv) PrivateKey {
	return PrivateKey{env: env, ptr: 0}
}

// ToString converts the PublicKey to its string representation using the linked Wasm environment. Returns the string or an error.
func (privateKey PrivateKey) ToString() (string, error) {
	if privateKey.ptr == 0 {
		return "", fmt.Errorf("token private_key not initialized")
	}

	return shared.AsString(privateKey.env, privateKey)
}

// FromString converts the string representation of a PrivateKey to a PrivateKey using the linked Wasm environment. Returns an error.
func (privateKey PrivateKey) FromString(env wasm.WasmEnv, data string) (PrivateKey, error) {

	returnArea, err := env.GetReturnArea()
	if err != nil {
		return PrivateKey{}, err
	}

	strPtr, err := env.WriteBytesToWasm([]byte(data))

	// Call: privatekey_fromString(out_ptr, str_ptr, str_len)
	_, err = env.Call("private_key_from_hex", returnArea, strPtr, uint64(len(data)))
	if err != nil {
		return PrivateKey{}, fmt.Errorf("privatekey_fromString failed: %w", err)
	}

	ptr, err := env.ResultPointer(returnArea)
	if err != nil {
		return PrivateKey{}, err
	}

	privateKey.ptr = ptr
	privateKey.env = env
	return privateKey, nil
}
