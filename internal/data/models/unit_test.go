package models_test

import (
	"testing"
	"time"

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

func TestIORateString(t *testing.T) {
	input := models.IORate(500)
	output := "500 B/s"
	if input.String() != output {
		t.Errorf("got %s, expected: %s", input, output)
	}
}

func TestBitRateString(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.BitRate
		output string
	}{
		{
			desc:   "if value is less than a kilobit (1000), String() should end with 'bps'",
			input:  models.BitRate(999.9),
			output: "999.9 bps",
		},
		{
			desc:   "if value is less than a megabit (1000000), String() should end with 'Kbps'",
			input:  models.BitRate(900000.0),
			output: "900.0 Kbps",
		},
		{
			desc:   "if value is less than a gigabit (1000000000), String() should end with 'Mbps'",
			input:  models.BitRate(900000000.0),
			output: "900.0 Mbps",
		},
		{
			desc:   "if value is more than a gigabit (1000000000), String() should end with 'Gbps'",
			input:  models.BitRate(9000000000.0),
			output: "9.0 Gbps",
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

func TestDistanceString(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.Distance
		output string
	}{
		{
			desc:   "if value is less than kilometer (1000), String() should end with 'm'",
			input:  models.Distance(999),
			output: "999.0 m",
		},
		{
			desc:   "if value is more than kilometer (1000), String() should end with 'Km'",
			input:  models.Distance(1000.1),
			output: "1.0 Km",
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

func TestDurationString(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.Duration
		output string
	}{
		{
			desc:   "if value is less than a microsecond, String() should end with 'ns'",
			input:  models.Duration(91 * time.Nanosecond),
			output: "91 ns",
		},
		{
			desc:   "if value is less than a millisecond, String() should end with 'us'",
			input:  models.Duration(87 * time.Microsecond),
			output: "87 us",
		},
		{
			desc:   "if value is less than a second, String() should end with 'ms'",
			input:  models.Duration(999 * time.Millisecond),
			output: "999 ms",
		},
		{
			desc:   "if value is less than a minute, String() should end with 's'",
			input:  models.Duration(50 * time.Second),
			output: "50 s",
		},
		{
			desc:   "if value is less than an hour, String() should end with 'min'",
			input:  models.Duration(50 * time.Minute),
			output: "50 min",
		},
		{
			desc:   "if value is more than an hour, String() should have h m s format",
			input:  models.Duration(72 * time.Minute),
			output: "1h 12m 0s",
		},
		{
			desc:   "if value is more than an hour, String() should have d h format",
			input:  models.Duration(1500 * time.Minute),
			output: "1d 1h",
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
