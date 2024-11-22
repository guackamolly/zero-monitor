package models_test

import (
	"testing"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

func TestJoinNetworkCodeExpiry(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.JoinNetworkCode
		output bool
	}{
		{
			desc:   "returns false if time now is before expire time",
			input:  models.JoinNetworkCode{Code: models.UUID(), ExpiresAt: time.Now().Add(time.Minute)},
			output: false,
		},
		{
			desc:   "returns true if time now is after expire time",
			input:  models.JoinNetworkCode{Code: models.UUID(), ExpiresAt: time.Now().Add(-time.Minute)},
			output: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if output := tC.input.Expired(); output != tC.output {
				t.Errorf("expected %v but got %v", tC.output, output)
			}
		})
	}
}
