package token_test

import (
	"biscuit-wasm-go/factory"
	"biscuit-wasm-go/keypair"
	tokenModule "biscuit-wasm-go/token"
	"fmt"
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

	fmt.Println(expectedBase64)

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

func TestBiscuit_FromBytes_FromFile(t *testing.T) {

	data, err := os.ReadFile("../samples/data/current/test001_basic.bc")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("------------------------Public key from string------------------------------------")

	publicKey, err := keypair.PublicKey{}.FromString(env, "acdd6d5b53bfee478bf689f8e012fe7988bf755e3d7c5152947abc149bc20189", keypair.Ed25519)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("------------------------From bytes------------------------------------")

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

	t.Log(b64)
	//
	//bytess, err := biscuit.ToBytes()
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//if len(bytess) != len(data) {
	//	t.Error("Expected the same length")
	//	return
	//}
	//for i := 0; i < len(bytess); i++ {
	//	if bytess[i] != data[i] {
	//		t.Error("Expected the same data")
	//	}
	//}

}

// TestBiscuit_FromBytes verifies the creation of a biscuit from bytes.
func TestBiscuit_FromBytes(t *testing.T) {

	code := "right(\"file1\", \"read\");\nright(\"file2\", \"read\");\nright(\"file1\", \"write\")\n"

	fmt.Println("------------------------Make Biscuit------------------------------------")
	token, err := factory.MakeBiscuit(env, code)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("------------------------To bytes------------------------------------")
	expectedBytes, err := token.Token.ToBytes()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("---------------------- To base64 ---------------------------------------")
	expectedBase64, err := token.Token.ToBase64()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("-----------------------------From bytes -----------------------------------")

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
