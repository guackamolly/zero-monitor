package models_test

import (
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/shirou/gopsutil/net"
)

func TestSpeedtestNextPhase(t *testing.T) {
	speedtest := models.NewSpeedtest("nk873-56-d2-355", "zero-monitor", "okla", "srv-01", 10)
	initialSpeedtest := speedtest
	nextPhase := func() models.Speedtest {
		speedtest = speedtest.NextPhase()
		return speedtest
	}

	testCases := []struct {
		desc   string
		input  models.Speedtest
		output models.SpeedtestPhase
	}{
		{
			desc:   "initial phase should be init",
			input:  initialSpeedtest,
			output: models.SpeedtestInit,
		},
		{
			desc:   "second phase should be latency test",
			input:  nextPhase(),
			output: models.SpeedtestLatency,
		},
		{
			desc:   "third phase should be download test",
			input:  nextPhase(),
			output: models.SpeedtestDownload,
		},
		{
			desc:   "fourth phase should be upload test",
			input:  nextPhase(),
			output: models.SpeedtestUpload,
		},
		{
			desc:   "last phase should be finish",
			input:  nextPhase(),
			output: models.SpeedtestFinish,
		},
		{
			desc:   "calling NextPhase after finish phase does nothing",
			input:  nextPhase(),
			output: models.SpeedtestFinish,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.input.Phase != tC.output {
				t.Errorf("expected %v, got %v", tC.output, tC.input.Phase)
			}
		})
	}
}

func TestSpeedtestFinished(t *testing.T) {
	speedtest := models.NewSpeedtest("nk873-56-d2-355", "zero-monitor", "okla", "srv-01", 10)
	initialSpeedtest := speedtest
	nextPhase := func() models.Speedtest {
		speedtest = speedtest.NextPhase()
		return speedtest
	}

	testCases := []struct {
		desc   string
		input  models.Speedtest
		output bool
	}{
		{
			desc:   "returns false if initial phase",
			input:  initialSpeedtest,
			output: false,
		},
		{
			desc:   "returns false if latency phase",
			input:  nextPhase(),
			output: false,
		},
		{
			desc:   "returns false if download phase",
			input:  nextPhase(),
			output: false,
		},
		{
			desc:   "returns false if upload phase",
			input:  nextPhase(),
			output: false,
		},
		{
			desc:   "returns true if finish phase",
			input:  nextPhase(),
			output: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.input.Finished() != tC.output {
				t.Errorf("expected %v, got %v", tC.output, tC.input.Finished())
			}
		})
	}
}

func TestSpeedtestFinishedLatency(t *testing.T) {
	speedtest := models.NewSpeedtest("nk873-56-d2-355", "zero-monitor", "okla", "srv-01", 10)
	initialSpeedtest := speedtest
	nextPhase := func() models.Speedtest {
		speedtest = speedtest.NextPhase()
		return speedtest
	}

	testCases := []struct {
		desc   string
		input  models.Speedtest
		output bool
	}{
		{
			desc:   "returns false if initial phase",
			input:  initialSpeedtest,
			output: false,
		},
		{
			desc:   "returns false if latency phase",
			input:  nextPhase(),
			output: false,
		},
		{
			desc:   "returns true if download phase",
			input:  nextPhase(),
			output: true,
		},
		{
			desc:   "returns true if upload phase",
			input:  nextPhase(),
			output: true,
		},
		{
			desc:   "returns true if finish phase",
			input:  nextPhase(),
			output: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.input.FinishedLatency() != tC.output {
				t.Errorf("expected %v, got %v", tC.output, tC.input.FinishedLatency())
			}
		})
	}
}

func TestSpeedtestFinishedDownload(t *testing.T) {
	speedtest := models.NewSpeedtest("nk873-56-d2-355", "zero-monitor", "okla", "srv-01", 10)
	initialSpeedtest := speedtest
	nextPhase := func() models.Speedtest {
		speedtest = speedtest.NextPhase()
		return speedtest
	}

	testCases := []struct {
		desc   string
		input  models.Speedtest
		output bool
	}{
		{
			desc:   "returns false if initial phase",
			input:  initialSpeedtest,
			output: false,
		},
		{
			desc:   "returns false if latency phase",
			input:  nextPhase(),
			output: false,
		},
		{
			desc:   "returns false if download phase",
			input:  nextPhase(),
			output: false,
		},
		{
			desc:   "returns true if upload phase",
			input:  nextPhase(),
			output: true,
		},
		{
			desc:   "returns true if finish phase",
			input:  nextPhase(),
			output: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.input.FinishedDownload() != tC.output {
				t.Errorf("expected %v, got %v", tC.output, tC.input.FinishedUpload())
			}
		})
	}
}

func TestSpeedtestFinishedUpload(t *testing.T) {
	speedtest := models.NewSpeedtest("nk873-56-d2-355", "zero-monitor", "okla", "srv-01", 10)
	initialSpeedtest := speedtest
	nextPhase := func() models.Speedtest {
		speedtest = speedtest.NextPhase()
		return speedtest
	}

	testCases := []struct {
		desc   string
		input  models.Speedtest
		output bool
	}{
		{
			desc:   "returns false if initial phase",
			input:  initialSpeedtest,
			output: false,
		},
		{
			desc:   "returns false if latency phase",
			input:  nextPhase(),
			output: false,
		},
		{
			desc:   "returns false if download phase",
			input:  nextPhase(),
			output: false,
		},
		{
			desc:   "returns false if upload phase",
			input:  nextPhase(),
			output: false,
		},
		{
			desc:   "returns true if finish phase",
			input:  nextPhase(),
			output: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.input.FinishedUpload() != tC.output {
				t.Errorf("expected %v, got %v", tC.output, tC.input.FinishedUpload())
			}
		})
	}
}

func TestConnectionTCP(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.Connection
		output bool
	}{
		{
			desc:   "returns true if connection kind is 1 (TCP)",
			input:  models.NewConnection(1, "none", net.Addr{}, net.Addr{}),
			output: true,
		},
		{
			desc:   "returns false if connection kind is not 1 (TCP)",
			input:  models.NewConnection(0, "none", net.Addr{}, net.Addr{}),
			output: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if out := tC.input.TCP(); out != tC.output {
				t.Errorf("expected %v, but got %v", tC.output, out)
			}
		})
	}
}

func TestConnectionUDP(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.Connection
		output bool
	}{
		{
			desc:   "returns true if connection kind is 2 (UDP)",
			input:  models.NewConnection(2, "none", net.Addr{}, net.Addr{}),
			output: true,
		},
		{
			desc:   "returns false if connection kind is not 2 (UDP)",
			input:  models.NewConnection(0, "none", net.Addr{}, net.Addr{}),
			output: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if out := tC.input.UDP(); out != tC.output {
				t.Errorf("expected %v, but got %v", tC.output, out)
			}
		})
	}
}

func TestConnectionExposed(t *testing.T) {
	testCases := []struct {
		desc   string
		input  models.Connection
		output bool
	}{
		{
			desc:   "returns true if connection local address IP is 0.0.0.0",
			input:  models.NewConnection(2, "none", net.Addr{IP: "0.0.0.0"}, net.Addr{}),
			output: true,
		},
		{
			desc:   "returns true if connection local address IP is ::",
			input:  models.NewConnection(2, "none", net.Addr{IP: "::"}, net.Addr{}),
			output: true,
		},
		{
			desc:   "returns false if connection local address IP is 127.0.0.1",
			input:  models.NewConnection(2, "none", net.Addr{IP: "127.0.0.1"}, net.Addr{}),
			output: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if out := tC.input.Exposed(); out != tC.output {
				t.Errorf("expected %v, but got %v", tC.output, out)
			}
		})
	}
}
