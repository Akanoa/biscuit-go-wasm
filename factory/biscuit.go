package factory

import (
	builderModule "biscuit-wasm-go/builder"
	keyPairModule "biscuit-wasm-go/keypair"
	"biscuit-wasm-go/token"
	"biscuit-wasm-go/wasm"
)

type MakeBiscuitReturn struct {
	Token      token.Biscuit
	Keypair    keyPairModule.KeyPair
	PrivateKey keyPairModule.PrivateKey
	PublicKey  keyPairModule.PublicKey
}

func MakeBiscuit(env wasm.WasmEnv, code string) (MakeBiscuitReturn, error) {
	keypair, err := keyPairModule.KeyPair{}.New(env, keyPairModule.Ed25519)
	if err != nil {
		return MakeBiscuitReturn{}, err
	}

	builder, err := builderModule.BiscuitBuilder{}.New(env)
	if err != nil {
		return MakeBiscuitReturn{}, err
	}

	err = builder.AddCode(code)
	if err != nil {
		return MakeBiscuitReturn{}, err
	}

	err = builder.SetRootKeyId(666)
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
