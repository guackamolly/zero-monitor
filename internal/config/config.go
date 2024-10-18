package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type Config struct {
	TrustedNetwork map[string]models.Node
}

// Loads a previously saved configuration file from user config directory.
// It uses [os.UserConfigDir] as the base directory and saves the config file
// under zero-monitor/config.json.
func Load() (Config, error) {
	cfg := Config{TrustedNetwork: map[string]models.Node{}}

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
