package mq_test

import (
	"slices"
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/mq"
)

func TestCryptoRegisterCipherKey(t *testing.T) {
	identity := []byte(models.UUID())

	testCases := []struct {
		desc  string
		key   []byte
		error bool
	}{
		{
			desc:  "returns error if key size is not 128/192/256 bits",
			key:   make([]byte, 4),
			error: true,
		},
		{
			desc:  "does not return error if key size is 128 bits",
			key:   make([]byte, 16),
			error: false,
		},
		{
			desc:  "does not return error if key size is 192 bits",
			key:   make([]byte, 24),
			error: false,
		},
		{
			desc:  "does not return error if key size is 256 bits",
			key:   make([]byte, 32),
			error: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if error := mq.RegisterCipherKey(identity, tC.key) != nil; error != tC.error {
				t.Errorf("expected %v but got %v", tC.error, error)
			}
		})
	}
}

func TestCryptoEncryptCipher(t *testing.T) {
	identity := []byte(models.UUID())
	data := []byte("zero-monitor")

	key, err := mq.GenerateCipherKey()
	if err != nil {
		t.Fatalf("didn't expect generate cipher key to fail, %v", err)
	}

	err = mq.RegisterCipherKey(identity, key)
	if err != nil {
		t.Fatalf("didn't expect register cipher key to fail, %v", err)
	}

	testCases := []struct {
		desc     string
		identity []byte
		error    bool
	}{
		{
			desc:     "returns error if identity is not associated to a cipher block",
			identity: []byte{},
			error:    true,
		},
		{
			desc:     "does not return error if identity is associated to a cipher block",
			identity: identity,
			error:    false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, _, err := mq.EncryptCipher(tC.identity, data)
			if error := err != nil; error != tC.error {
				t.Errorf("expected %v but got %v", tC.error, error)
			}
		})
	}
}

func TestCryptoDecryptCipher(t *testing.T) {
	identity := []byte(models.UUID())
	data := []byte("zero-monitor")

	key, err := mq.GenerateCipherKey()
	if err != nil {
		t.Fatalf("didn't expect generate cipher key to fail, %v", err)
	}

	err = mq.RegisterCipherKey(identity, key)
	if err != nil {
		t.Fatalf("didn't expect register cipher key to fail, %v", err)
	}

	cipher, nonce, err := mq.EncryptCipher(identity, data)
	if err != nil {
		t.Fatalf("didn't expect encrypt cipher to fail, %v", err)
	}

	testCases := []struct {
		desc     string
		identity []byte
		nonce    []byte
		error    bool
	}{
		{
			desc:     "returns error if identity is not associated to a cipher block",
			identity: []byte{},
			nonce:    nonce,
			error:    true,
		},
		{
			desc:     "returns error if nonce does not match nonce used to encrypt",
			identity: identity,
			nonce:    []byte{},
			error:    true,
		},
		{
			desc:     "does not return error if identity is associated to a cipher block and nonce matches encrypted message",
			identity: identity,
			nonce:    nonce,
			error:    false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := mq.DecryptCipher(tC.identity, cipher, tC.nonce)
			if error := err != nil; error != tC.error {
				t.Errorf("expected %v but got %v", tC.error, error)
			}
		})
	}
}

func TestEncryptAndDecryptCipher(t *testing.T) {
	identity := []byte(models.UUID())
	data := []byte("zero-monitor")

	key, err := mq.GenerateCipherKey()
	if err != nil {
		t.Fatalf("didn't expect generate cipher key to fail, %v", err)
	}

	err = mq.RegisterCipherKey(identity, key)
	if err != nil {
		t.Fatalf("didn't expect register cipher key to fail, %v", err)
	}

	cipher, nonce, err := mq.EncryptCipher(identity, data)
	if err != nil {
		t.Fatalf("didn't expect encrypt cipher to fail, %v", err)
	}

	plain, err := mq.DecryptCipher(identity, cipher, nonce)
	if err != nil {
		t.Fatalf("didn't expect decrypt cipher to fail, %v", err)
	}

	if !slices.Equal(plain, data) {
		t.Errorf("expected decrypt to return %v, but got %v", data, plain)
	}
}
