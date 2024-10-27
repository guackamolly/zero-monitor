package mq

import (
	"encoding/gob"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type JoinNetworkRequest struct {
	Node models.Node
}

type JoinNetworkResponse struct {
	StatsPoll time.Duration
}

type UpdateNodeStatsRequest struct {
	Node models.Node
}

type UpdateNodeStatsPollDurationRequest struct {
	Duration time.Duration
}

type NodeConnectionsResponse struct {
	Connections []models.Connection
}

type NodeProcessesResponse struct {
	Processes []models.Process
}

func init() {
	gob.Register(JoinNetworkRequest{})
	gob.Register(JoinNetworkResponse{})
	gob.Register(UpdateNodeStatsRequest{})
	gob.Register(UpdateNodeStatsPollDurationRequest{})
	gob.Register(NodeConnectionsResponse{})
	gob.Register(NodeProcessesResponse{})

	gob.Register(models.Node{})

}
