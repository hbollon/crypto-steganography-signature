package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
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
// Export bool param save the private key to a file if true.
// It returns an Signer and an Unsigner.
// It also returns an error if the key pair, the Signer or the Unsigner could not be generated.
func InitializeKeyPair(bitsSize int, export bool) (signer Signer, unsigner Unsigner, err error) {
	// Generate keypair
	privateKey, err := rsa.GenerateKey(rand.Reader, bitsSize)
	if err != nil {
		return nil, nil, err
	}
	if export {
		exportKeyPairToPem(privateKey, &privateKey.PublicKey)
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

// InitializeKeyPairFromFiles loads existing RSA keypair from given pem files.
// It initialize and returns an Signer and an Unsigner.
// It also returns an error if the key pair, the Signer or the Unsigner could not be generated.
func InitializeKeyPairFromFiles(privateKeyPath, publicKeyPath string) (signer Signer, unsigner Unsigner, err error) {
	// Load private key
	rawPrivateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return
	}
	block, _ := pem.Decode(rawPrivateKey)
	if block == nil {
		return nil, nil, fmt.Errorf("no private key found")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	// Load public key
	rawPublicKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return
	}
	block, _ = pem.Decode(rawPublicKey)
	if block == nil {
		return nil, nil, fmt.Errorf("no public key found")
	}
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return
	}

	// Create signer
	signer, err = newSignerFromKey(privateKey)
	if err != nil {
		return nil, nil, err
	}

	// Create unsigner
	unsigner, err = newUnsignerFromKey(publicKey)
	if err != nil {
		return nil, nil, err
	}

	return
}

func exportKeyPairToPem(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) error {
	// Export private key to file
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)
	f_priv, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer f_priv.Close()
	_, err = f_priv.Write(privateKeyPem)
	if err != nil {
		return err
	}

	// Export public key to file
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)
	f_pub, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer f_pub.Close()
	_, err = f_pub.Write(publicKeyPem)
	if err != nil {
		return err
	}

	return nil
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
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)
	return rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA256, d, sig)
}
