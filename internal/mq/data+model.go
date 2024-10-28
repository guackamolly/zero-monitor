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

type KillNodeProcessRequest struct {
	PID int32
}

type KillNodeProcessResponse struct {
	Processes []models.Process
}

type OPError struct {
	Err string
}

func (e *OPError) Error() string {
	return e.Err
}

func init() {
	gob.Register(JoinNetworkRequest{})
	gob.Register(JoinNetworkResponse{})
	gob.Register(UpdateNodeStatsRequest{})
	gob.Register(UpdateNodeStatsPollDurationRequest{})
	gob.Register(NodeConnectionsResponse{})
	gob.Register(NodeProcessesResponse{})
	gob.Register(KillNodeProcessRequest{})
	gob.Register(KillNodeProcessResponse{})

	gob.Register(models.Node{})
	gob.Register(OPError{})

}
