package models_test

import (
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

func TestSpeedtestNextPhase(t *testing.T) {
	speedtest := models.NewSpeedtest("zero-monitor", "srv-01")
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
	speedtest := models.NewSpeedtest("zero-monitor", "srv-01")
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
	speedtest := models.NewSpeedtest("zero-monitor", "srv-01")
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
	speedtest := models.NewSpeedtest("zero-monitor", "srv-01")
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
	speedtest := models.NewSpeedtest("zero-monitor", "srv-01")
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
