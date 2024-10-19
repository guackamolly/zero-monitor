package service

import (
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

type NodeSchedulerService struct{}

// Creates a new service that schedules routines that interact with
// the nodes network.
// All routines are initialized when calling this function.
//
// This is a dumb service and shouldn't interact with no injected instance. Instead
// instances should be retrieved by callbacks.
func NewNodeSchedulerService(
	saveCfg func() error,
	updateTrustedNetwork func([]models.Node),
	networkStream func() chan ([]models.Node),
) *NodeSchedulerService {
	// schedule goroutine that reacts to network changes and updates config
	go func() {
		for s := range networkStream() {
			logging.LogInfo("nodes network changed, updating trusted network")
			updateTrustedNetwork(s)
		}
	}()

	// schedule goroutine that save config every 5 minutes
	go func() {
		time.Sleep(5 * time.Minute)

		logging.LogInfo("trying to save config file")
		err := saveCfg()

		if err != nil {
			logging.LogError("coudln't save config file, %v", err)
		}
	}()

	return &NodeSchedulerService{}
}
