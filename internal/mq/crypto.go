package mq

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// Global map of cipher blocks used to encrypt/decrypt sensitive
// messages between pub/sub nodes.
//
// Key is the pub/sub identity
var blocks = map[string]cipher.Block{}

// Registers a key and creates new cipher block.
func RegisterCipherKey(
	identity []byte,
	key []byte,
) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	blocks[string(identity)] = block
	return nil
}

// Encrypts data using the cipher block associated to the pub/sub identity.
// Returns encrypted data followed by the nonce used to authenticate the data.
func EncryptCipher(
	identity []byte,
	data []byte,
) ([]byte, []byte, error) {
	block, ok := blocks[string(identity)]
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
	block, ok := blocks[string(identity)]
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

// Generates a new 256 bit cipher key.
func GenerateCipherKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate key, %v", err)
	}

	return key, nil
}
