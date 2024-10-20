package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

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
	for _, n := range nodes {
		c.TrustedNetwork[n.ID] = n
	}
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

func configJsonPath() (string, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return "", errors.Join(errors.New("couldn't lookup user config dir"), err)
	}

	p := filepath.Join(d, "zero-monitor", "config.json")
	err = os.MkdirAll(filepath.Dir(p), 0700)
	if err != nil {
		return p, errors.Join(errors.New("couldn't stat config file dir"), err)
	}

	return p, nil
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
