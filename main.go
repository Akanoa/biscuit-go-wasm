package main

import (
	keypairModule "biscuit-wasm-go/crypto/keypair"
	"biscuit-wasm-go/wasm"
	"fmt"
	"log/slog"
	"os"
)

func createkeypair(env wasm.WasmEnv, algorithm keypairModule.SignatureAlgorithm) (keypairModule.KeyPair, error) {
	keypair := keypairModule.KeyPair{}
	if err := keypair.New(env, algorithm); err != nil {
		slog.Error(err.Error())
		return keypair, err
	}

	privateKey, err := keypair.GetPrivateKey()
	if err != nil {
		slog.Error(err.Error())
		return keypair, err
	}

	privateKeyString, err := privateKey.ToString()
	if err != nil {
		slog.Error(err.Error())
		return keypair, err
	}
	fmt.Printf("PrivateKeyString %s\n", privateKeyString)

	return keypair, nil
}

func main() {

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	slog.SetDefault(logger)

	env, err := wasm.InitWasm()
	if err != nil {
		panic(err)
	}

	//keypair1, err := createkeypair(env, keypairModule.Ed25519)
	_, err = createkeypair(env, keypairModule.Ed25519)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	//privateKey1, err := keypair1.GetPrivateKey()
	//if err != nil {
	//	println("keypair1.GetPrivateKey error:", err.Error())
	//	return
	//}
	//
	//privateKey1String, err := privateKey1.ToString()
	//if err != nil {
	//	println("privateKey1.ToString error:", err.Error())
	//}
	//fmt.Println("From Keypair", privateKey1String)

	//keypair2 := keypairModule.Invoke(module, ctx)
	//err = keypair2.FromPrivateKey(privateKey1)
	//if err != nil {
	//	println("keypair2.FromPrivateKey error:", err.Error())
	//	return
	//}
	//
	//privateKey2, err := keypair2.GetPrivateKey()
	//if err != nil {
	//	println("keypair2.GetPrivateKey error:", err.Error())
	//	return
	//}
	//
	//privateKey2String, err := privateKey2.ToString()
	//if err != nil {
	//	println("privateKey2.ToString error:", err.Error())
	//}
	//
	//fmt.Println("From Keypair", privateKey2String)

	privateKey3 := keypairModule.InvokePrivateKey(env)
	err = privateKey3.FromString("ed25519-private/eacbce4ed1a4132e1c667ebe5f730f493197fd3def32027a87ea2233d5b55abb")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	privateKey3String, err := privateKey3.ToString()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	fmt.Println("From PrivateKey", privateKey3String)

}
