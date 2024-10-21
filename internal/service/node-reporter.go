package service

import (
	"log"
	"time"

	"github.com/guackamolly/zero-monitor/internal"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Service for reporting node information to master.
type NodeReporterService struct {
	initial models.Node
	stream  chan (models.Node)
	system  repositories.SystemRepository
}

func NewNodeReporterService(
	system repositories.SystemRepository,
) *NodeReporterService {
	s := &NodeReporterService{
		system: system,
		stream: make(chan models.Node),
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
func (s NodeReporterService) Start(pollDuration time.Duration) chan (models.Node) {
	go func() {
		node := s.initial
		for {
			stats, err := s.system.Stats()
			if err != nil {
				log.Printf("failed to fetch system statistics, %v\n", err)
			} else {
				node = node.WithUpdatedStats(stats)
				s.stream <- node
			}

			logging.LogInfo("sleeping for %s until polling new node stats", pollDuration)
			time.Sleep(pollDuration)
		}
	}()

	return s.stream
}

// tries to fetch system info, and if it fails, sleeps for 2 seconds until trying again.
// blocking call.
func (s NodeReporterService) systemInfo() models.Info {
	for {
		sinfo, err := s.system.Info()
		if err == nil {
			return sinfo
		}

		log.Printf("failed to fetch system information, %v\n", err)
		time.Sleep(time.Second * 2)
	}
}
