package keypair

import (
	"strings"
	"testing"
)

func TestKeyPairNonDeterministic(t *testing.T) {
	keyPair1, err := KeyPair{}.New(env, Ed25519)
	if err != nil {
		t.Error(err)
		return
	}
	keyPair2, err := KeyPair{}.New(env, Ed25519)
	if err != nil {
		t.Error(err)
		return
	}

	privateKey1, err := keyPair1.GetPrivateKey()
	if err != nil {
		t.Error(err)
		return
	}
	privateKey2, err := keyPair2.GetPrivateKey()
	if err != nil {
		t.Error(err)
		return
	}

	privateKey1HexString, err := privateKey1.ToString()
	if err != nil {
		t.Error(err)
		return
	}

	privateKey2HexString, err := privateKey2.ToString()
	if err != nil {
		t.Error(err)
		return
	}

	if privateKey1HexString == privateKey2HexString {
		t.Errorf("Two new private keys mustn't be the same !")
	}
}

// TestKeyPairSignatureAlgorithmEd25519 verifies the creation of a keypair using the Ed25519 algorithm and its private key format.
func TestKeyPairSignatureAlgorithmEd25519(t *testing.T) {
	keyPair, err := KeyPair{}.New(env, Ed25519)
	if err != nil {
		t.Error(err)
		return
	}
	privateKey, err := keyPair.GetPrivateKey()
	if err != nil {
		t.Error(err)
	}
	privateKeyString, err := privateKey.ToString()
	if err != nil {
		t.Error(err)
	}
	if strings.HasPrefix(privateKeyString, "ed25519-private/") == false {
		t.Errorf("Expected private key to start with ed25519-private/")
	}
}

// TestKeyPairSignatureAlgorithmSecp256r1 verifies the creation of a keypair using the Secp256r1 algorithm and its private key format.
func TestKeyPairSignatureAlgorithmSecp256r1(t *testing.T) {
	keyPair, err := KeyPair{}.New(env, Secp256r1)
	if err != nil {
		t.Error(err)
		return
	}
	privateKey, err := keyPair.GetPrivateKey()
	if err != nil {
		t.Error(err)
	}
	privateKeyString, err := privateKey.ToString()
	if err != nil {
		t.Error(err)
	}
	if strings.HasPrefix(privateKeyString, "secp256r1-private/") == false {
		t.Errorf("Expected private key to start with secp256r1-private/; found: %s", privateKeyString)
	}
}

// TestKeyPair_FromPrivateKeyEd25519 verifies the creation of a keypair from a private key.
func TestKeyPair_FromPrivateKeyEd22519(t *testing.T) {

	testimonialKeypair, err := KeyPair{}.New(env, Ed25519)
	if err != nil {
		t.Error(err)
		return
	}

	testimonialPrivateKey, err := testimonialKeypair.GetPrivateKey()
	if err != nil {
		t.Error(err)
		return
	}
	testimonialPrivateKeyString, err := testimonialPrivateKey.ToString()
	if err != nil {
		t.Error(err)
		return
	}

	keypair, err := KeyPair{}.FromPrivateKey(env, testimonialPrivateKey)
	if err != nil {
		t.Error(err)
	}

	privateKey, err := keypair.GetPrivateKey()
	if err != nil {
		t.Error(err)
	}
	privateKeyString, err := privateKey.ToString()
	if err != nil {
		t.Error(err)
	}
	if testimonialPrivateKeyString != privateKeyString {
		t.Errorf("Expected %s, got %s", testimonialPrivateKeyString, privateKeyString)
	}

}

// TestKeyPair_FromPrivateKeyEd25519 verifies the creation of a keypair from a private key.
func TestKeyPair_FromPrivateKeySecp256r1(t *testing.T) {

	testimonialKeypair, err := KeyPair{}.New(env, Secp256r1)
	if err != nil {
		t.Error(err)
		return
	}

	testimonialPrivateKey, err := testimonialKeypair.GetPrivateKey()
	if err != nil {
		t.Error(err)
		return
	}
	testimonialPrivateKeyString, err := testimonialPrivateKey.ToString()
	if err != nil {
		t.Error(err)
		return
	}

	keypair, err := KeyPair{}.FromPrivateKey(env, testimonialPrivateKey)
	if err != nil {
		t.Error(err)
	}

	privateKey, err := keypair.GetPrivateKey()
	if err != nil {
		t.Error(err)
	}
	privateKeyString, err := privateKey.ToString()
	if err != nil {
		t.Error(err)
	}
	if testimonialPrivateKeyString != privateKeyString {
		t.Errorf("Expected %s, got %s", testimonialPrivateKeyString, privateKeyString)
	}

}
