package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

// Signer can create signatures that verify against a public key.
type Signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
}

// Unsigner can verify signatures made with Signer.
type Unsigner interface {
	// Verify verifies signature validity with existing data. This method
	// will apply the existing signature specified to the data.
	Verify(message []byte, sig []byte) error
}

// InitializeKeyPair generates a new RSA keypair of the given bit size using the random source random.
// It returns an Signer and an Unsigner.
// It also returns an error if the key pair, the Signer or the Unsigner could not be generated.
func InitializeKeyPair(bitsSize int) (signer Signer, unsigner Unsigner, err error) {
	// Generate keypair
	privateKey, err := rsa.GenerateKey(rand.Reader, bitsSize)
	if err != nil {
		return nil, nil, err
	}

	// Create signer
	signer, err = newSignerFromKey(privateKey)
	if err != nil {
		return nil, nil, err
	}

	// Create unsigner
	unsigner, err = newUnsignerFromKey(privateKey.Public())
	if err != nil {
		return nil, nil, err
	}

	return
}

func newSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

func newUnsignerFromKey(k interface{}) (Unsigner, error) {
	var sshKey Unsigner
	switch t := k.(type) {
	case *rsa.PublicKey:
		sshKey = &rsaPublicKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

type rsaPublicKey struct {
	*rsa.PublicKey
}

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}

// Verify verifies signature validity with existing data
func (r *rsaPublicKey) Verify(message []byte, sig []byte) error {
	return rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA256, message, sig)
}
