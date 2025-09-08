package factory

import (
	builderModule "biscuit-wasm-go/builder"
	keyPairModule "biscuit-wasm-go/keypair"
	"biscuit-wasm-go/token"
	"biscuit-wasm-go/wasm"
	"fmt"
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

	fmt.Println("--------------------------------------------")

	err = builder.AddCode(code)
	if err != nil {
		return MakeBiscuitReturn{}, err
	}
	fmt.Println("--------------------------------------------")

	privateKey, err := keypair.GetPrivateKey()
	if err != nil {
		return MakeBiscuitReturn{}, err
	}

	x, err := privateKey.ToString()
	fmt.Println("Private:", x)

	publicKey, err := keypair.GetPublicKey()
	if err != nil {
		return MakeBiscuitReturn{}, err
	}

	x, err = publicKey.ToString()
	fmt.Println("Public:", x)

	biscuit, err := builder.Build(privateKey)
	if err != nil {
		return MakeBiscuitReturn{}, err
	}

	fmt.Println("--------------------------------------------")

	return MakeBiscuitReturn{
		Token:      biscuit,
		Keypair:    keypair,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}
