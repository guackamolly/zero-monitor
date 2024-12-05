package mq

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

// Global map of blocks used to encrypt/decrypt sensitive
// messages between pub/sub nodes.
//
// Key is the pub/sub identity
var cipherBlocks = map[string]cipher.Block{}

// The public/private key block used to encrypt/decrypt cipher keys during
// the key exchange between nodes.
var blk *pem.Block

// Loads the public/private key block to be used on encryption/decryption.
func LoadAsymmetricBlock(
	keyfilepath string,
) error {
	f, err := os.ReadFile(keyfilepath)
	if err != nil {
		return err
	}

	blk, _ = pem.Decode(f)
	return nil
}

// Registers a key and creates new cipher block.
func RegisterCipherKey(
	identity []byte,
	key []byte,
) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	cipherBlocks[string(identity)] = block
	return nil
}

// Encrypts data using the cipher block associated to the pub/sub identity.
// Returns encrypted data followed by the nonce used to authenticate the data.
func EncryptCipher(
	identity []byte,
	data []byte,
) ([]byte, []byte, error) {
	block, ok := cipherBlocks[string(identity)]
	if !ok {
		return nil, nil, fmt.Errorf("no encryption block associated with identity: %x", identity)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce, %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to construct aes-gcm cipher, %v", err)
	}

	bs := gcm.Seal(nil, nonce, data, nil)
	return bs, nonce, nil
}

// Decrypts encrypted data using the cipher block associated to the pub/sub identity.
func DecryptCipher(
	identity []byte,
	data []byte,
	nonce []byte,
) (plain []byte, err error) {
	block, ok := cipherBlocks[string(identity)]
	if !ok {
		return nil, fmt.Errorf("no encryption block associated with identity: %x", identity)
	}

	gcm, xerr := cipher.NewGCM(block)
	if xerr != nil {
		return nil, fmt.Errorf("failed to construct aes-gcm cipher, %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("open panic, %v", r)
		}
	}()
	return gcm.Open(nil, nonce, data, nil)
}

// Encrypts data using the public key block.
// Block must have been loaded before by calling [LoadAsymmetricBlock(true)].
func EncryptAsymmetric(
	data []byte,
) ([]byte, error) {
	key, err := x509.ParsePKIXPublicKey(blk.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, key.(*rsa.PublicKey), data)
}

// Encrypts data using the private key block.
// Block must have been loaded before by calling [LoadAsymmetricBlock(false)].
func DecryptAsymmetric(
	data []byte,
) ([]byte, error) {
	key, err := x509.ParsePKCS1PrivateKey(blk.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, key, data)
}

// Generates a new 256 bit cipher key.
func GenerateCipherKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate key, %v", err)
	}

	return key, nil
}

// TODO: derive from private key instead
func DerivePublicKey() ([]byte, error) {
	key, err := x509.ParsePKCS1PrivateKey(blk.Bytes)
	if err != nil {
		return nil, err
	}

	// Extract the public key from the private key
	publicKey := &key.PublicKey

	// Marshal the public key into DER format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	// Encode the public key into a PEM block
	publicKeyBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return publicKeyBlock, nil
}
