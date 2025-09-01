package keypair

// Placeholder test file previously had incomplete references that broke `go test`.
// Keeping the package testable without undefined symbols.

import (
	"biscuit-wasm-go/wasm"
	"log"
	"testing"
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

// Ensure that the private key can be converted to and from a string.
func TestPrivateKey_FromString(t *testing.T) {
	testimonialPrivateKeyHexString := "ed25519-private/ece6f591177602076ea2ec48ac73063a038959bfa1f2e6794a0d076febad34b1"

	privateKey := PrivateKey{}.New(env)
	err := privateKey.FromString(testimonialPrivateKeyHexString)
	if err != nil {
		t.Error(err)
	}

	privateKeyString, err := privateKey.ToString()
	if err != nil {
		t.Error(err)
		return
	}
	if testimonialPrivateKeyHexString != privateKeyString {
		t.Errorf("Expected %s, got %s", testimonialPrivateKeyHexString, privateKeyString)
	}
}
