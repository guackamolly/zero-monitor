package service

import (
	"time"

	"github.com/guackamolly/zero-monitor/internal"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Service for reporting node information to master.
type NodeReporterService struct {
	initial           models.Node
	statsPollDuration time.Duration
	system            repositories.SystemRepository
}

func NewNodeReporterService(
	system repositories.SystemRepository,
) *NodeReporterService {
	s := &NodeReporterService{
		system: system,
	}

	id := internal.MachineId
	info := s.systemInfo()
	s.initial = models.NewNodeWithoutStats(id, info)

	return s
}

func (s NodeReporterService) Initial() models.Node {
	return s.initial
}

// Starts reporting node information through a channel. The channel is unbuffered.
func (s *NodeReporterService) Start(pollDuration time.Duration) chan (models.Node) {
	stream := make(chan (models.Node))
	s.statsPollDuration = pollDuration

	go func() {
		node := s.initial
		for {
			pollDuration := s.statsPollDuration
			stats, err := s.system.Stats()
			if err != nil {
				logging.LogError("failed to fetch system statistics, %v", err)
			} else {
				node = node.WithUpdatedStats(stats)
				stream <- node
			}

			logging.LogInfo("sleeping for %s until polling new node stats", pollDuration)
			time.Sleep(pollDuration)
		}
	}()

	return stream
}

func (s *NodeReporterService) Update(
	d time.Duration,
) {
	s.statsPollDuration = d
}

// tries to fetch system info, and if it fails, sleeps for 2 seconds until trying again.
// blocking call.
func (s NodeReporterService) systemInfo() models.Info {
	for {
		sinfo, err := s.system.Info()
		if err == nil {
			return sinfo
		}

		logging.LogError("failed to fetch system information, %v", err)
		time.Sleep(time.Second * 2)
	}
}
