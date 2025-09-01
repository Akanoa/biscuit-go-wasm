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
	function, err := env.GetFunction("keypair_new")
	if err != nil {
		return keypair, err
	}

	result, err := env.Call(function, uint64(signatureAlgorithm))
	if err != nil {
		return keypair, fmt.Errorf("keypair_new failed: %w", err)
	}

	if len(result) == 0 {
		return keypair, fmt.Errorf("no result returned from keypair_new")
	}

	keypair.ptr = result[0]
	keypair.env = env

	return keypair, nil
}

// GetPublicKey returns the public key of the keypair.
func (keypair *KeyPair) GetPublicKey() (PublicKey, error) {

	if keypair.ptr == 0 {
		slog.Error("keypair not initialized")
		return PublicKey{}, fmt.Errorf("keypair not initialized")
	}

	function, err := keypair.env.GetFunction("keypair_getPublicKey")
	if err != nil {
		slog.Error("exported function 'keypair_getPublicKey' not found")
		return PublicKey{}, err
	}

	result, err := keypair.env.Call(function, keypair.ptr)
	if err != nil {
		slog.Error("keypair_getPublicKey failed", slog.Any("err", err))
		return PublicKey{}, err
	}

	return PublicKey{
		ptr: result[0],
		env: keypair.env,
	}, nil
}

// GetPrivateKey returns the private key of the keypair.
func (keypair *KeyPair) GetPrivateKey() (PrivateKey, error) {

	if keypair.ptr == 0 {
		return PrivateKey{}, fmt.Errorf("keypair not initialized")
	}

	function, err := keypair.env.GetFunction("keypair_getPrivateKey")
	if err != nil {
		slog.Error("exported function 'keypair_getPrivateKey' not found")
		return PrivateKey{}, err
	}

	result, err := keypair.env.Call(function, keypair.ptr)
	if err != nil {
		slog.Error("keypair_getPrivateKey failed", slog.Any("err", err))
		return PrivateKey{}, err
	}

	return PrivateKey{
		ptr: result[0],
		env: keypair.env,
	}, nil
}

// FromPrivateKey creates a new keypair from the specified private key.
func (keypair KeyPair) FromPrivateKey(env wasm.WasmEnv, privateKey PrivateKey) (KeyPair, error) {

	function, err := env.GetFunction("keypair_fromPrivateKey")
	if err != nil {
		slog.Error("exported function 'keypair_fromPrivateKey' not found")
		return keypair, err
	}

	result, err := env.Call(function, privateKey.ptr)

	if err != nil {
		slog.Error("keypair_fromPrivateKey failed", slog.Any("err", err))
		return keypair, err
	}

	if len(result) == 0 {
		return keypair, fmt.Errorf("no result returned from keypair_fromPrivateKey")
	}

	keypair.ptr = result[0]
	keypair.env = env

	return keypair, nil
}
