package service

import (
	"time"

	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

type NodeSchedulerService struct{}

// Creates a new service that schedules routines that interact with
// the nodes network.
// All routines are initialized when calling this function.
func NewNodeSchedulerService(
	cfg *config.Config,
	networkStream func() chan ([]models.Node),
) *NodeSchedulerService {
	// schedule goroutine that reacts to network changes and updates config
	go func() {
		for s := range networkStream() {
			logging.LogInfo("nodes network changed, updating memory config")

			for _, n := range s {
				cfg.TrustedNetwork[n.ID] = n
			}
		}
	}()

	// schedule goroutine that save config every 5 minutes
	go func() {
		time.Sleep(5 * time.Minute)

		logging.LogInfo("trying to save config file")
		err := config.Save(*cfg)

		if err != nil {
			logging.LogError("coudln't save config file, %v", err)
		}
	}()

	return &NodeSchedulerService{}
}
