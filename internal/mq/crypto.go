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
	"path/filepath"

	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

const (
	mqTransportPrivateKeyFileEnvKey = "mq_transport_pem_key"
	mqTransportPublicKeyFileEnvKey  = "mq_transport_pub_key"
)

var (
	mqTransportPrivateKeyFile = os.Getenv(mqTransportPrivateKeyFileEnvKey)
	mqTransportPublicKeyFile  = os.Getenv(mqTransportPublicKeyFileEnvKey)
)

// Global map of blocks used to encrypt/decrypt sensitive
// messages between pub/sub nodes.
//
// Key is the pub/sub identity
var cipherBlocks = map[string]cipher.Block{}

// The public/private key block used to encrypt/decrypt cipher keys during
// the key exchange between nodes.
var pemBlock *pem.Block

func init() {
	peml := len(mqTransportPrivateKeyFile)
	publ := len(mqTransportPublicKeyFile)

	if peml > 0 && publ > 0 {
		return
	}

	d, err := config.Dir()
	if err != nil {
		logging.LogWarning("couldn't lookup pem/pub key files to encrypt message queue messages. either communication with master node fail OR it won't be encrypted")
		return
	}

	if peml == 0 {
		mqTransportPrivateKeyFile = filepath.Join(d, "mq.pem")
	}

	if publ == 0 {
		mqTransportPublicKeyFile = filepath.Join(d, "mq.pub")
	}
}

// Loads the public/private key block to be used on encryption/decryption.
func LoadAsymmetricBlock(
	pub bool,
) error {
	p := mqTransportPrivateKeyFile
	if pub {
		p = mqTransportPublicKeyFile
	}

	f, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	pemBlock, _ = pem.Decode(f)
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
	key, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
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
	key, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
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
