package token_test

import (
	tokenModule "biscuit-wasm-go/token"
	"testing"
)

func TestBiscuit_FromBase64(t *testing.T) {
	code := "user(1)"

	token, err := MakeBisuit(code)
	if err != nil {
		t.Error(err)
	}

	expectedBase64, err := token.Token.ToBase64()
	if err != nil {
		t.Error(err)
		return
	}

	biscuit, err := tokenModule.Biscuit{}.FromBase64(env, expectedBase64, token.PublicKey)
	if err != nil {
		t.Error(err)
		return
	}

	base64, err := biscuit.ToBase64()
	if err != nil {
		t.Error(err)
		return
	}

	if expectedBase64 != base64 {
		t.Errorf("Expected %s, got %s", expectedBase64, base64)
	}

}
