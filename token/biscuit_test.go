package token_test

import (
	"biscuit-wasm-go/factory"
	tokenModule "biscuit-wasm-go/token"
	"testing"
)

// TestBiscuit_FromBase64 verifies the creation of a biscuit from base64.
func TestBiscuit_FromBase64(t *testing.T) {
	code := "user(1)"

	token, err := factory.MakeBiscuit(env, code)
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

// TestBiscuit_FromBytes verifies the creation of a biscuit from bytes.
func TestBiscuit_FromBytes(t *testing.T) {
	code := "user(1)"

	token, err := factory.MakeBiscuit(env, code)
	if err != nil {
		t.Error(err)
	}

	expectedBase64, err := token.Token.ToBase64()
	if err != nil {
		t.Error(err)
		return
	}

	expectedBytes, err := token.Token.ToBytes()
	if err != nil {
		t.Error(err)
	}

	biscuit, err := tokenModule.Biscuit{}.FromBytes(env, expectedBytes, token.PublicKey)
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
