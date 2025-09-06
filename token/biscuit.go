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

	ret, err := env.Call(function, strPtr, uint64(len(biscuitBase64)), publicKey.Ptr())
	if err != nil {
		return Biscuit{}, fmt.Errorf("biscuit_fromBase64 failed: %w", err)
	}

	valuePtr, err := env.GetPointee(ret)

	if err != nil {
		return Biscuit{}, err
	}

	biscuit.env = env
	biscuit.ptr = valuePtr
	return biscuit, nil
}

// ToBytes returns the byte slice representation of the Biscuit.
func (biscuit Biscuit) ToBytes() ([]byte, error) {
	if biscuit.ptr == 0 {
		return nil, fmt.Errorf("token not initialized")
	}

	function, err := biscuit.env.GetFunction("biscuit_toBytes")
	if err != nil {
		return nil, err
	}

	returnPtr, err := biscuit.env.GetReturnArea()
	defer biscuit.env.Free(returnPtr, wasm.ReturnAreaSize)

	_, err = biscuit.env.Call(function, biscuit.ptr)
	if err != nil {
		return nil, err
	}

	return biscuit.env.GetBytesValueFromPointer(returnPtr)
}

// PadToMultipleOf4 renvoie une nouvelle slice dont la taille est un multiple de 4
func PadToMultipleOf4(data []byte) []byte {
	n := len(data)
	remainder := n % 8
	if remainder == 0 {
		return data // Déjà multiple de 4
	}
	padding := 4 - remainder
	padded := make([]byte, n+padding)
	copy(padded, data)
	// Les octets restants sont automatiquement 0
	return padded
}

// FromBytes creates a new Biscuit from a byte slice.
func (biscuit Biscuit) FromBytes(env wasm.WasmEnv, bytes []byte, publicKey keypair.PublicKey) (Biscuit, error) {

	function, err := env.GetFunction("biscuit_fromBytes")
	if err != nil {
		return Biscuit{}, err
	}

	dataPtr, err := env.WriteBytes(bytes)
	fmt.Printf("----ptr %x data: %v\n", dataPtr, bytes)
	if err != nil {
		return Biscuit{}, err
	}
	defer func() {
		fmt.Println("free bytes")
		err := env.Free(dataPtr, uint64(len(bytes)))
		if err != nil {
			fmt.Printf("cannot free bytes: %v\n", err)
		}

	}()

	fmt.Println("BEFORE")
	ret, err := env.Call(function, dataPtr, uint64(len(bytes)), publicKey.Ptr())

	fmt.Println("ret", ret)
	fmt.Println("err", err)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return Biscuit{}, err
	}
	fmt.Println("AFTER")

	if err != nil {
		return Biscuit{}, err
	}
	return Biscuit{}.New(env, ret[0]), nil

}
