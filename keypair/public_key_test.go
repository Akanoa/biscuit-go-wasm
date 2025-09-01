package keypair

import (
	"strings"
	"testing"
)

// Ensure that the public key can be converted to and from a string.
// Using Ed25519 signature algorithm
func TestPublicKey_FromStringEd25519(t *testing.T) {
	keypair, err := KeyPair{}.New(env, Ed25519)
	if err != nil {
		t.Error(err)
		return
	}

	expectedPublicKey, err := keypair.GetPublicKey()
	if err != nil {
		t.Error(err)
		return
	}

	expectedPublicKeyString, err := expectedPublicKey.ToString()
	if err != nil {
		t.Error(err)
		return
	}

	_, expectedPublicKeyStringWithoutAlgorithm, found := strings.Cut(expectedPublicKeyString, "/")
	if found == false {
		t.Errorf("Expected algorithm to be found in %s", expectedPublicKeyString)
		return
	}

	publicKey, err := PublicKey{}.FromString(env, expectedPublicKeyStringWithoutAlgorithm, Ed25519)
	if err != nil {
		t.Error(err)
		return
	}

	publicKeyString, err := publicKey.ToString()
	if err != nil {
		t.Error(err)
	}

	if expectedPublicKeyString != publicKeyString {
		t.Errorf("Expected %s, got %s", expectedPublicKeyString, publicKeyString)
	}
}

// Ensure that the public key can be converted to and from a string.
// Using Secp256r1 signature algorithm
func TestPublicKey_FromStringSecp256r1(t *testing.T) {
	keypair, err := KeyPair{}.New(env, Secp256r1)
	if err != nil {
		t.Error(err)
		return
	}

	expectedPublicKey, err := keypair.GetPublicKey()
	if err != nil {
		t.Error(err)
		return
	}

	expectedPublicKeyString, err := expectedPublicKey.ToString()
	if err != nil {
		t.Error(err)
		return
	}

	_, expectedPublicKeyStringWithoutAlgorithm, found := strings.Cut(expectedPublicKeyString, "/")
	if found == false {
		t.Errorf("Expected algorithm to be found in %s", expectedPublicKeyString)
		return
	}

	publicKey, err := PublicKey{}.FromString(env, expectedPublicKeyStringWithoutAlgorithm, Secp256r1)
	if err != nil {
		t.Error(err)
		return
	}

	publicKeyString, err := publicKey.ToString()
	if err != nil {
		t.Error(err)
	}

	if expectedPublicKeyString != publicKeyString {
		t.Errorf("Expected %s, got %s", expectedPublicKeyString, publicKeyString)
	}
}
