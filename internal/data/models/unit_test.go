package models_test

import (
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

func TestPercentString(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.Percent
		output string
	}{
		{
			desc:   "if value is smaller than 0, String() should return '-'",
			input:  models.Percent(-1),
			output: "-",
		},
		{
			desc:   "if value is 0, String() should return '0.00%'",
			input:  models.Percent(0),
			output: "0.00%",
		},
		{
			desc:   "if value is greater than 0, String() should return percentage with two decimal cases",
			input:  models.Percent(2.556),
			output: "2.56%",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.input.String() != tC.output {
				t.Errorf("got %s, expected: %s", tC.input, tC.output)
			}
		})
	}
}

func TestMemoryString(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.Memory
		output string
	}{
		{
			desc:   "if value is smaller than a kilobyte (1024), String() should end with 'B'",
			input:  models.Memory(1023),
			output: "1023 B",
		},
		{
			desc:   "if value is greater than a kilobyte, but smaller than a megabyte (1048576), String() should end with 'KB'",
			input:  models.Memory(1048575),
			output: "1023 KB",
		},
		{
			desc:   "if value is greater than a megabyte, but smaller than a gigabyte (1073741824), String() should end with 'MB'",
			input:  models.Memory(1073741823),
			output: "1023 MB",
		},
		{
			desc:   "if value is greater than a gigabyte, but smaller than a terabyte (1099511627776), String() should end with 'GB'",
			input:  models.Memory(1099511627775),
			output: "1023 GB",
		},
		{
			desc:   "if value is equal or greater than a terabyte, String() should end with 'TB'",
			input:  models.Memory(1099511627776),
			output: "1 TB",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.input.String() != tC.output {
				t.Errorf("got %s, expected: %s", tC.input, tC.output)
			}
		})
	}
}

func TestCelsiusString(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.Celsius
		output string
	}{
		{
			desc:   "if value cannot be rounded to an higher value, String() returns <integer> ºC",
			input:  models.Celsius(50.5),
			output: "50 ºC",
		},
		{
			desc:   "if value can be rounded to an higher value, String() returns <integer + 1> ºC",
			input:  models.Celsius(50.51),
			output: "51 ºC",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.input.String() != tC.output {
				t.Errorf("got %s, expected: %s", tC.input, tC.output)
			}
		})
	}
}
