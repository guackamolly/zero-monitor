package service

import (
	"log"
	"time"

	"github.com/guackamolly/zero-monitor/internal"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
)

// Service for reporting node information to master.
type NodeReporterService struct {
	stream chan (models.Node)
	system repositories.SystemRepository
}

func NewNodeReporterService(
	system repositories.SystemRepository,
) *NodeReporterService {
	return &NodeReporterService{
		system: system,
		stream: make(chan models.Node),
	}
}

// Starts reporting node information through a channel. The channel is unbuffered.
func (s NodeReporterService) Start() chan (models.Node) {
	go func() {
		id := internal.MachineId
		info := s.systemInfo()
		node := models.NewNodeWithoutStats(id, info)

		for {
			stats, err := s.system.Stats()
			if err != nil {
				log.Printf("failed to fetch system statistics, %v\n", err)
				time.Sleep(time.Second * 5)
				continue
			}

			node = node.WithUpdatedStats(stats)
			s.stream <- node
			time.Sleep(time.Second * 5)
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
