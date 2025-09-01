package builder

import (
	keyPairModule "biscuit-wasm-go/keypair"
	"testing"
)

func TestBiscuitBuilder_Build(t *testing.T) {
	keypair, err := keyPairModule.KeyPair{}.New(env, keyPairModule.Ed25519)
	if err != nil {
		t.Error(err)
		return
	}
	privatekey, err := keypair.GetPrivateKey()
	if err != nil {
		t.Error(err)
	}

	biscuitBuilder, err := BiscuitBuilder{}.New(env)
	if err != nil {
		t.Error(err)
		return
	}

	code := "user(1)"

	err = biscuitBuilder.AddCode(code)
	if err != nil {
		t.Error(err)
		return
	}

	err = biscuitBuilder.SetRootKeyId(666)
	if err != nil {
		t.Error(err)
		return
	}

	str, err := biscuitBuilder.ToString()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(str)

	biscuit, err := biscuitBuilder.Build(privatekey)
	if err != nil {
		t.Error(err)
		return
	}

	bisuitBase64, err := biscuit.ToBase64()
	if err != nil {
		t.Error(err)
	}
	t.Log(bisuitBase64)

}
