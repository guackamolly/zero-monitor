package models

import (
	"testing"
	"time"
)

func TestJoinNetworkCodeExpiry(t *testing.T) {
	testCases := []struct {
		desc   string
		input  JoinNetworkCode
		output bool
	}{
		{
			desc:   "returns false if time now is before expire time",
			input:  JoinNetworkCode{Code: UUID(), ExpiresAt: time.Now().Add(time.Minute)},
			output: false,
		},
		{
			desc:   "returns true if time now is after expire time",
			input:  JoinNetworkCode{Code: UUID(), ExpiresAt: time.Now().Add(-time.Minute)},
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
