package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

var machineID string

type configurableValue struct {
	Value   int
	Default int
	Min     int
	Max     int
}

type Config struct {
	TrustedNetwork      map[string]models.Node
	NodeStatsPolling    configurableValue
	NodeLastSeenTimeout configurableValue
	AutoSavePeriod      configurableValue
}

func (c *Config) UpdateTrustedNetwork(
	nodes []models.Node,
) {
	network := map[string]models.Node{}
	for _, n := range nodes {
		network[n.ID] = n
	}

	c.TrustedNetwork = network
}

func (c *Config) UpdateConfigurableValues(
	nodeStatsPolling int,
	nodeLastSeenTimeout int,
	autoSavePeriod int,
) {
	c.NodeStatsPolling.Value = nodeStatsPolling
	c.NodeLastSeenTimeout.Value = nodeLastSeenTimeout
	c.AutoSavePeriod.Value = autoSavePeriod
}

func (cv configurableValue) Duration() time.Duration {
	return time.Duration(cv.Value) * time.Second
}

// Loads a previously saved configuration file from user config directory.
// It uses [os.UserConfigDir] as the base directory and saves the config file
// under zero-monitor/config.json.
func Load() (Config, error) {
	cfg := defaultConfig()

	p, err := configJsonPath()
	if err != nil {
		return cfg, err
	}

	bs, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return cfg, err
	}

	if err != nil {
		return cfg, errors.Join(errors.New("couldn't read config file"), err)
	}

	err = json.Unmarshal(bs, &cfg)
	return cfg, err
}

func Save(cfg Config) error {
	p, err := configJsonPath()
	if err != nil {
		return err
	}

	bs, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(p, bs, 0666)
}

// Returns the path to the config directory used by the program.
// If error is not nil, then an error occurred when creating the directory.
func Dir() (string, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return "", errors.Join(errors.New("couldn't lookup user config dir"), err)
	}

	p := filepath.Join(d, "zero-monitor")
	err = os.MkdirAll(filepath.Dir(p), 0700)
	if err != nil {
		return p, errors.Join(errors.New("couldn't stat config file dir"), err)
	}

	return p, nil
}

// Returns the machine id used to identify node agents.
func MachineID() string {
	if machineID != "" {
		return machineID
	}

	d, err := Dir()
	if err != nil {
		logging.LogFatal("couldn't stat config directory, which is required for extracting the machine unique id!, %v", err)
	}

	p := filepath.Join(d, "machine-id")
	_, err = os.Stat(p)
	if os.IsNotExist(err) {
		machineID = models.UUID()
		err = os.WriteFile(p, []byte(machineID), os.ModePerm)
	}

	if err != nil {
		logging.LogFatal("couldn't stat machine id file, which is required to identify node agents!, %v", err)
	}

	bs, err := os.ReadFile(p)
	if err != nil {
		logging.LogFatal("couldn't read machine id file, which is required to identify node agents!, %v", err)
	}

	machineID = string(bs)
	return machineID
}

func configJsonPath() (string, error) {
	d, err := Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(d, "config.json"), nil
}

func defaultConfig() Config {
	return Config{
		TrustedNetwork: map[string]models.Node{},
		NodeStatsPolling: configurableValue{
			Value:   5,
			Default: 5,
			Min:     1,
			Max:     60 * 10,
		},
		NodeLastSeenTimeout: configurableValue{
			Value:   10,
			Default: 10,
			Min:     5,
			Max:     60 * 10,
		},
		AutoSavePeriod: configurableValue{
			Value:   600,
			Default: 600,
			Min:     60,
			Max:     600 * 6 * 24,
		},
	}
}
