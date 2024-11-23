package service_test

import (
	"testing"

	"github.com/guackamolly/zero-monitor/internal/service"
)

func TestNetworkServiceCodeReturnsLastCodeCreated(t *testing.T) {
	s := service.NewNetworkService(nil)

	// create new code
	c1 := s.Code()
	c2 := s.Code()

	if c1 != c2 {
		t.Errorf("expected %v, but got %v", c1, c2)
	}
}

func TestNetworkServiceValid(t *testing.T) {
	s := service.NewNetworkService(nil)

	testCases := []struct {
		desc   string
		input  string
		output bool
	}{
		{
			desc:   "returns false if no code has been created yet",
			input:  "my.code",
			output: false,
		},
		{
			desc:   "returns true if code matches current code",
			input:  func() string { return s.Code().Code }(),
			output: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if output := s.Valid(tC.input); output != tC.output {
				t.Errorf("expected %v but got %v", output, tC.output)
			}
		})
	}
}
