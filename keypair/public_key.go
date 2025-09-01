package keypair

import (
	"biscuit-wasm-go/shared"
	"biscuit-wasm-go/wasm"
	"fmt"
	"log/slog"
)

// PublicKey represents a WASM public key.
type PublicKey struct {
	env wasm.WasmEnv
	ptr uint64
}

// ToStringWasmFunction returns the name of the WASM function that converts the PublicKey to its string representation.
func (publicKey PublicKey) ToStringWasmFunction() string {
	return "publickey_toString"
}

// Ptr returns the pointer to the underlying WASM object.
func (public_key PublicKey) Ptr() uint64 {
	return public_key.ptr
}

// FromString initializes a PublicKey from a string and a given SignatureAlgorithm using a Wasm environment.
// Returns the initialized PublicKey or an error in case of failure.
func (publicKey PublicKey) FromString(env wasm.WasmEnv, data string, algorithm SignatureAlgorithm) (PublicKey, error) {
	function, err := env.GetFunction("publickey_fromString")
	if err != nil {
		slog.Error("exported function 'publickey_fromString' not found")
		return publicKey, err
	}

	retPtr, err := env.GetReturnArea()
	if err != nil {
		return publicKey, fmt.Errorf("malloc for return area failed: %w", err)
	}
	defer env.Free(retPtr, wasm.ReturnAreaSize)

	strPtr, err := env.WriteString(data)
	if err != nil {
		return publicKey, fmt.Errorf("cannot write string to wasm memory: %w", err)
	}
	defer env.Free(strPtr, uint64(len(data)))

	_, err = env.Call(function, retPtr, strPtr, uint64(len(data)), uint64(algorithm))

	// Read result triple
	valuePtr, err := env.GetPointee(retPtr)
	if err != nil {
		return publicKey, err
	}

	publicKey.ptr = valuePtr
	publicKey.env = env
	return publicKey, nil
}

// ToString converts the PublicKey to its string representation using the linked Wasm environment. Returns the string or an error.
func (publick_key PublicKey) ToString() (string, error) {
	if publick_key.ptr == 0 {
		return "", fmt.Errorf("biscuit publick_key not initialized")
	}

	return shared.AsString(publick_key.env, publick_key)
}
