package service

import (
	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Service that acts as a facade for configuration requests.
type MasterConfigurationService struct {
	cfg *config.Config
}

func NewMasterConfigurationService(
	cfg *config.Config,
) *MasterConfigurationService {
	return &MasterConfigurationService{
		cfg: cfg,
	}
}

func (s MasterConfigurationService) Current() config.Config {
	return *s.cfg
}

// Updates the trusted network present in the current configuration instance.
func (s MasterConfigurationService) UpdateTrustedNetwork(
	nodes []models.Node,
) {
	for _, n := range nodes {
		s.cfg.TrustedNetwork[n.ID] = n
	}
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
}

func (s MasterConfigurationService) Save() error {
	err := config.Save(*s.cfg)
	if err != nil {
		logging.LogError("failed to save config file, %v", err)
	}

	return err
}
