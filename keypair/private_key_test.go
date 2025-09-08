package keypair

import (
	"testing"
)

// Ensure that the private key can be converted to and from a string.
func TestPrivateKey_FromString(t *testing.T) {

	testimonialPrivateKeyHexString := "ed25519-private/ece6f591177602076ea2ec48ac73063a038959bfa1f2e6794a0d076febad34b1"

	privateKey, err := PrivateKey{}.FromString(env, testimonialPrivateKeyHexString)
	if err != nil {
		t.Error(err)
		return
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
