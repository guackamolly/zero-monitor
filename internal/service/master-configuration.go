package service

import (
	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Service that acts as a facade for configuration requests.
type MasterConfigurationService struct {
	cfg    *config.Config
	stream chan (config.Config)
}

func NewMasterConfigurationService(
	cfg *config.Config,
) *MasterConfigurationService {
	return &MasterConfigurationService{
		cfg:    cfg,
		stream: make(chan config.Config),
	}
}

func (s MasterConfigurationService) Current() config.Config {
	return *s.cfg
}

func (s MasterConfigurationService) Stream() chan (config.Config) {
	return s.stream
}

// Updates the trusted network present in the current configuration instance.
func (s MasterConfigurationService) UpdateTrustedNetwork(
	nodes []models.Node,
) {
	s.cfg.UpdateTrustedNetwork(nodes)

	s.updateStream()
}

// Updates all configurable values present in the current configuration instance.
func (s MasterConfigurationService) UpdateConfigurable(
	nodeStatsPolling int,
	nodeLastSeenTimeout int,
	autoSavePeriod int,
) {
	s.cfg.UpdateConfigurableValues(
		nodeStatsPolling,
		nodeLastSeenTimeout,
		autoSavePeriod,
	)

	s.updateStream()
}

func (s MasterConfigurationService) Save() error {
	err := config.Save(*s.cfg)
	if err != nil {
		logging.LogError("failed to save config file, %v", err)
	}

	return err
}

func (s MasterConfigurationService) updateStream() {
	go func() {
		s.stream <- *s.cfg
	}()
}
