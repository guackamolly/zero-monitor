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
	node              models.Node
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
	s.node = models.NewNodeWithoutStats(id, info)

	return s
}

func (s NodeReporterService) Node() models.Node {
	return s.node
}

// Starts reporting node information through a channel. The channel is unbuffered.
func (s *NodeReporterService) Start(pollDuration time.Duration) chan (models.Node) {
	stream := make(chan (models.Node))
	s.statsPollDuration = pollDuration

	go func() {
		for {
			pollDuration := s.statsPollDuration
			stats, err := s.system.Stats()
			if err != nil {
				logging.LogError("failed to fetch system statistics, %v", err)
			} else {
				s.node = s.node.WithUpdatedStats(stats)
				stream <- s.node
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

func (s NodeReporterService) Connections() ([]models.Connection, error) {
	return s.system.Conns()
}

func (s NodeReporterService) Processes() ([]models.Process, error) {
	return s.system.Procs()
}

func (s NodeReporterService) KillProcess(pid int32) error {
	return s.system.KillProc(pid)
}

// tries to fetch system info, and if it fails, sleeps for 2 seconds until trying again.
// blocking call.
func (s NodeReporterService) systemInfo() models.MachineInfo {
	for {
		sinfo, err := s.system.Info()
		if err == nil {
			return sinfo
		}

		logging.LogError("failed to fetch system information, %v", err)
		time.Sleep(time.Second * 2)
	}
}
