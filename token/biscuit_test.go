package token_test

import (
	"biscuit-wasm-go/factory"
	"biscuit-wasm-go/keypair"
	tokenModule "biscuit-wasm-go/token"
	"encoding/base64"
	"os"
	"testing"
)

// TestBiscuit_FromBase64 verifies the creation of a biscuit from base64.
func TestBiscuit_FromBase64(t *testing.T) {
	code := "user(1)"

	token, err := factory.MakeBiscuit(env, code)
	if err != nil {
		t.Error(err)
		return
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

	b64, err := biscuit.ToBase64()
	if err != nil {
		t.Error(err)
		return
	}

	if expectedBase64 != b64 {
		t.Errorf("Expected %s, got %s", expectedBase64, b64)
	}

}

func TestBiscuit_FromBytes_FromFile(t *testing.T) {

	data, err := os.ReadFile("../samples/data/current/test001_basic.bc")
	if err != nil {
		t.Error(err)
		return
	}

	publicKey, err := keypair.PublicKey{}.FromString(env, "1055c750b1a1505937af1537c626ba3263995c33a64758aaafb1275b0312e284", keypair.Ed25519)
	if err != nil {
		t.Error(err)
		return
	}

	biscuit, err := tokenModule.Biscuit{}.FromBytes(env, data, publicKey)
	if err != nil {
		t.Error("Expected error", err)
		return
	}

	b64, err := biscuit.ToBase64()
	if err != nil {
		t.Error(err)
		return
	}

	// Decode the b64 string coming from the biscuit encoded
	bytes, err := base64.URLEncoding.DecodeString(b64)
	if err != nil {
		t.Error(err)
		return
	}

	// Compare the bytes between the file contents and the decoded biscuit
	for i := 0; i < len(data); i++ {
		if data[i] != bytes[i] {
			t.Errorf("Expected %d, got %d", data[i], bytes[i])
		}
	}

}

// TestBiscuit_FromBytes verifies the creation of a biscuit from bytes.
func TestBiscuit_FromBytes(t *testing.T) {

	code := "right(\"file1\", \"read\");\nright(\"file2\", \"read\");\nright(\"file1\", \"write\")\n"

	token, err := factory.MakeBiscuit(env, code)
	if err != nil {
		t.Error(err)
		return
	}

	expectedBytes, err := token.Token.ToBytes()
	if err != nil {
		t.Error(err)
		return
	}

	expectedBase64, err := token.Token.ToBase64()
	if err != nil {
		t.Error(err)
		return
	}

	biscuit, err := tokenModule.Biscuit{}.FromBytes(env, expectedBytes, token.PublicKey)
	if err != nil {
		t.Error(err)
		return
	}

	b64, err := biscuit.ToBase64()
	if err != nil {
		t.Error(err)
		return
	}

	if expectedBase64 != b64 {
		t.Errorf("Expected %s, got %s", expectedBase64, b64)
	}

}
