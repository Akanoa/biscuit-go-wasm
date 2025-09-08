// Package keypair provides high-level wrappers for key pair creation and access
// backed by the Biscuit WebAssembly module.
package keypair

import (
	"biscuit-wasm-go/wasm"
	"fmt"
	"log/slog"
)

type SignatureAlgorithm int

const (
	Ed25519   SignatureAlgorithm = iota
	Secp256r1                    = iota
)

type KeyPair struct {
	env wasm.WasmEnv
	ptr uint64
}

// New creates a new keypair using the specified signature algorithm.
func (keypair KeyPair) New(env wasm.WasmEnv, signatureAlgorithm SignatureAlgorithm) (KeyPair, error) {

	returnArea, err := env.GetReturnArea()
	if err != nil {
		return keypair, err
	}

	_, err = env.Call("keypair_new", returnArea, uint64(signatureAlgorithm))
	if err != nil {
		return keypair, fmt.Errorf("keypair_new failed: %w", err)
	}

	ptr, err := env.ResultPointer(returnArea)
	if err != nil {
		return keypair, err
	}

	keypair.ptr = ptr
	keypair.env = env

	return keypair, nil
}

// GetPublicKey returns the public key of the keypair.
func (keypair *KeyPair) GetPublicKey() (PublicKey, error) {

	if keypair.ptr == 0 {
		return PublicKey{}, fmt.Errorf("keypair not initialized")
	}

	returnArea, err := keypair.env.GetReturnArea()
	if err != nil {
		return PublicKey{}, err
	}

	_, err = keypair.env.Call("keypair_public_key", returnArea, keypair.ptr)
	if err != nil {
		slog.Error("keypair_getPublicKey failed", slog.Any("err", err))
		return PublicKey{}, err
	}

	ptr, err := keypair.env.ResultPointer(returnArea)
	if err != nil {
		return PublicKey{}, err
	}

	return PublicKey{
		ptr: ptr,
		env: keypair.env,
	}, nil
}

// GetPrivateKey returns the private key of the keypair.
func (keypair *KeyPair) GetPrivateKey() (PrivateKey, error) {

	if keypair.ptr == 0 {
		return PrivateKey{}, fmt.Errorf("keypair not initialized")
	}

	returnArea, err := keypair.env.GetReturnArea()
	if err != nil {
		return PrivateKey{}, err
	}

	fmt.Println("keypair_private_key")
	fmt.Println(keypair.ptr)
	_, err = keypair.env.Call("keypair_private_key", returnArea, keypair.ptr)
	if err != nil {
		slog.Error("keypair_getPublicKey failed", slog.Any("err", err))
		return PrivateKey{}, err
	}
	fmt.Println("keypair_private_key done")

	ptr, err := keypair.env.ResultPointer(returnArea)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		ptr: ptr,
		env: keypair.env,
	}, nil
}

// FromPrivateKey creates a new keypair from the specified private key.
func (keypair KeyPair) FromPrivateKey(env wasm.WasmEnv, privateKey PrivateKey) (KeyPair, error) {

	returnArea, err := env.GetReturnArea()
	if err != nil {
		return KeyPair{}, err
	}

	_, err = env.Call("keypair_from_private_key", returnArea, privateKey.Ptr())
	if err != nil {
		slog.Error("keypair_fromPrivatekey failed", slog.Any("err", err))
		return KeyPair{}, err
	}

	ptr, err := env.ResultPointer(returnArea)
	if err != nil {
		return KeyPair{}, err
	}

	keypair.ptr = ptr
	keypair.env = env

	return keypair, nil
}
