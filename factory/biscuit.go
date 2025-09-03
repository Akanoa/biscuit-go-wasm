package factory

import (
	"biscuit-wasm-go/builder"
	"biscuit-wasm-go/keypair"
	"biscuit-wasm-go/token"
	"biscuit-wasm-go/wasm"
)

type MakeBiscuitReturn struct {
	Token      token.Biscuit
	Keypair    keypair.KeyPair
	PrivateKey keypair.PrivateKey
	PublicKey  keypair.PublicKey
}

func MakeBiscuit(env wasm.WasmEnv, code string) (MakeBiscuitReturn, error) {
	keypair, err := keypair.KeyPair{}.New(env, keypair.Ed25519)
	if err != nil {
		return MakeBiscuitReturn{}, err
	}

	builder, err := builder.BiscuitBuilder{}.New(env)
	if err != nil {
		return MakeBiscuitReturn{}, err
	}

	err = builder.AddCode(code)
	if err != nil {
		return MakeBiscuitReturn{}, err
	}

	privateKey, err := keypair.GetPrivateKey()
	if err != nil {
		return MakeBiscuitReturn{}, err
	}

	publicKey, err := keypair.GetPublicKey()
	if err != nil {
		return MakeBiscuitReturn{}, err
	}

	biscuit, err := builder.Build(privateKey)
	if err != nil {
		return MakeBiscuitReturn{}, err
	}

	return MakeBiscuitReturn{
		Token:      biscuit,
		Keypair:    keypair,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}
