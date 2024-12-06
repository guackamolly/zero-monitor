package models_test

import (
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

func TestUserIsAdmin(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.User
		output bool
	}{
		{
			desc:   "returns true if user has admin role",
			input:  models.User{Role: models.AdminRole},
			output: true,
		},
		{
			desc:   "returns false if user has guest role",
			input:  models.User{Role: models.GuestRole},
			output: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if output := tC.input.IsAdmin(); output != tC.output {
				t.Errorf("expected %v but got %v", tC.output, output)
			}
		})
	}
}
