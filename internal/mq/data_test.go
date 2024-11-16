package mq_test

import (
	"reflect"
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/mq"
)

func TestMsgEncryptAttachesNonceToMetadataField(t *testing.T) {
	identity := []byte(models.UUID())
	m := mq.Msg{Identity: identity}

	key, err := mq.GenerateCipherKey()
	if err != nil {
		t.Fatalf("didn't expect generate cipher key to fail, %v", err)
	}

	err = mq.RegisterCipherKey(identity, key)
	if err != nil {
		t.Fatalf("didn't expect register cipher key to fail, %v", err)
	}

	m, err = m.Encrypt()
	if err != nil {
		t.Fatalf("didnt't expect encrypt to fail, %v", err)
	}

	if _, ok := m.Metadata.([]byte); !ok {
		t.Errorf("expected metadata field to contain nonce, but got: %v", m.Metadata)
	}
}

func TestMsgEncryptAttachesEncryptMsgToDataField(t *testing.T) {
	identity := []byte(models.UUID())
	m := mq.Msg{Identity: identity}

	key, err := mq.GenerateCipherKey()
	if err != nil {
		t.Fatalf("didn't expect generate cipher key to fail, %v", err)
	}

	err = mq.RegisterCipherKey(identity, key)
	if err != nil {
		t.Fatalf("didn't expect register cipher key to fail, %v", err)
	}

	m, err = m.Encrypt()
	if err != nil {
		t.Fatalf("didnt't expect encrypt to fail, %v", err)
	}

	if _, ok := m.Data.([]byte); !ok {
		t.Errorf("expected data field to contain encrypted message, but got: %v", m.Data)
	}
}

func TestMsgDecrypt(t *testing.T) {
	testCases := []struct {
		desc  string
		input mq.Msg
		error bool
	}{
		{
			desc:  "returns error if metadata field does not contain nonce",
			input: mq.Msg{Metadata: "not nonce", Data: make([]byte, 64)},
			error: true,
		},
		{
			desc:  "returns error if data field does not contain encrypted message",
			input: mq.Msg{Metadata: make([]byte, 12), Data: "not encrypted message"},
			error: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := tC.input.Decrypt()
			if error := err != nil; error != tC.error {
				t.Errorf("expected %v but got %v", tC.error, error)
			}
		})
	}
}

func TestMsgDecryptReturnsMsgBeforeEncryption(t *testing.T) {
	identity := []byte(models.UUID())
	m := mq.Msg{Identity: identity}

	key, err := mq.GenerateCipherKey()
	if err != nil {
		t.Fatalf("didn't expect generate cipher key to fail, %v", err)
	}

	err = mq.RegisterCipherKey(identity, key)
	if err != nil {
		t.Fatalf("didn't expect register cipher key to fail, %v", err)
	}

	em, err := m.Encrypt()
	if err != nil {
		t.Fatalf("didnt't expect encrypt to fail, %v", err)
	}

	dm, err := em.Decrypt()
	if err != nil {
		t.Fatalf("didnt't expect decrypt to fail, %v", err)
	}

	if !reflect.DeepEqual(m, dm) {
		t.Errorf("expected %v to be equal to %v", m, dm)
	}
}
