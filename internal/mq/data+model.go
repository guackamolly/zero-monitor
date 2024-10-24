package mq

import (
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
