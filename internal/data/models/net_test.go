package models_test

import (
	"net"
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

func TestIPString(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.IP
		output string
	}{
		{
			desc:   "if ip is less than 4 characters, String() should return '-'",
			input:  models.IP{255, 255, 255},
			output: "-",
		},
		{
			desc:   "if ip is 4 characters, String() should return IP as IPv4",
			input:  models.IP{255, 255, 255, 255},
			output: "255.255.255.255",
		},
		{
			desc:   "if ip is 16 characters, String() should return IP as IPv6",
			input:  models.IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			output: "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
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

func TestConnectionKindString(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.ConnectionKind
		output string
	}{
		{
			desc:   "if kind is 1, String() returns 'TCP'",
			input:  models.ConnectionKind(1),
			output: "TCP",
		},
		{
			desc:   "if kind is 2, String() returns 'UDP'",
			input:  models.ConnectionKind(2),
			output: "UDP",
		},
		{
			desc:   "if kind is neither 1 or 2, String() returns 'UNKNOWN'",
			input:  models.ConnectionKind(0),
			output: "UNKNOWN",
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

func TestNewNetAddress(t *testing.T) {
	testCases := []struct {
		desc  string
		input net.Addr
		error bool
	}{
		{
			desc:  "does not return error if tcp net addr",
			input: &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 65535},
			error: false,
		},
		{
			desc:  "does not return error if udp net addr",
			input: &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 65535},
			error: false,
		},
		{
			desc:  "returns error if any other addr",
			input: &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)},
			error: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := models.NewNetAddress(tC.input)
			if error := err != nil; error != tC.error {
				t.Errorf("expected %v but got %v", tC.error, error)
			}
		})
	}
}
