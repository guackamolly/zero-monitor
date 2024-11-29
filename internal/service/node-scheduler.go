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
//
// This is a dumb service and shouldn't interact with no injected instance. Instead
// instances should be retrieved by callbacks.
func NewNodeSchedulerService(
	cfg func() config.Config,
	cfgStream func() chan (config.Config),
	saveCfg func() error,
	updateTrustedNetwork func([]models.Node),
	updateNetwork func(models.Node) error,
	network func() []models.Node,
	networkStream func() chan ([]models.Node),
) *NodeSchedulerService {
	// schedule goroutine that reacts to network changes and updates config
	go func() {
		for s := range networkStream() {
			logging.LogDebug("nodes network changed, updating trusted network")
			updateTrustedNetwork(s)
		}
	}()

	// schedule goroutine that save config every 5 minutes
	go func() {
		for {
			autoSaveDuration := cfg().AutoSavePeriod.Duration()
			logging.LogDebug("sleeping for %s before saving config file", autoSaveDuration)
			time.Sleep(autoSaveDuration)

			logging.LogDebug("trying to save config file")
			err := saveCfg()

			if err != nil {
				logging.LogError("couldn't save config file, %v", err)
			}
		}
	}()

	// schedule goroutine that checks for any networ node that have gone missing
	go func() {
		for {
			t := time.Now()
			n := network()
			lastSeenTimeout := cfg().NodeLastSeenTimeout.Duration()
			for _, n := range n {
				if n.LastSeen.Sub(t).Abs() < lastSeenTimeout {
					continue
				}

				if !n.Online {
					continue
				}

				err := updateNetwork(n.SetOffline())
				if err != nil {
					logging.LogError("very strange error when notifying network that node is offline, %v", err)
				}
			}

			logging.LogDebug("sleeping for %s before checking for missing nodes", lastSeenTimeout)
			time.Sleep(lastSeenTimeout)
		}
	}()

	return &NodeSchedulerService{}
}
