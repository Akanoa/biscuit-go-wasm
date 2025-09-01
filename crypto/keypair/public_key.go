package keypair

import (
	"biscuit-wasm-go/wasm"
	"fmt"
	"log/slog"
)

type PublicKey struct {
	env wasm.WasmEnv
	ptr uint64
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
func (publicKey PublicKey) ToString() (string, error) {
	if publicKey.ptr == 0 {
		return "", fmt.Errorf("public key not initialized")
	}

	function, err := publicKey.env.GetFunction("public_key_ToString")
	if err != nil {
		return "", err
	}

	resultPtr, err := publicKey.env.GetStringArea()
	if err != nil {
		return "", err
	}
	defer publicKey.env.Free(resultPtr, wasm.StringAreaSize)

	_, err = publicKey.env.Call(function, resultPtr, publicKey.ptr)
	if err != nil {
		slog.Error("public_key_ToString failed", slog.Any("err", err))
		return "", err
	}

	return publicKey.env.GetStringValueFromPointer(resultPtr)
}
