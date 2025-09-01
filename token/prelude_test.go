package token_test

import (
	builderModule "biscuit-wasm-go/builder"
	keypairModule "biscuit-wasm-go/keypair"
	tokenModule "biscuit-wasm-go/token"
	"biscuit-wasm-go/wasm"
	"log"
)

var env wasm.WasmEnv

func init() {
	log.Println("setupSuite")
	result, err := wasm.InitWasm()
	if err != nil {
		log.Fatal(err)
		return
	}
	env = result
}

type MakeBiscuitReturn struct {
	Token      tokenModule.Biscuit
	Keypair    keypairModule.KeyPair
	PrivateKey keypairModule.PrivateKey
	PublicKey  keypairModule.PublicKey
}

func MakeBisuit(code string) (MakeBiscuitReturn, error) {
	keypair, err := keypairModule.KeyPair{}.New(env, keypairModule.Ed25519)
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
