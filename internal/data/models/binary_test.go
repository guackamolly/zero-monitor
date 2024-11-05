package models_test

import (
	"errors"
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

func TestEncodeErrorsIfGobDoesntRecognizeType(t *testing.T) {
	d := errors.New("unsupported type")
	_, err := models.Encode(d)
	if err == nil {
		t.Error("expected Encode() to error on unsupported type")
	}
}

func TestDecodeErrorsIfGobDoesntRecognizeType(t *testing.T) {
	bs := []byte("not an error")
	_, err := models.Decode[error](bs)
	if err == nil {
		t.Error("expected Decode() to error on unsupported type")
	}
}
