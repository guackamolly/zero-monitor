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

type AuthenticateNetworkRequest struct {
	Node       models.Node
	InviteCode string
}

type AuthenticateNetworkResponse struct{}
type RequiresAuthenticationResponse struct{}

type UpdateNodeStatsRequest struct {
	Stats models.Stats
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

type NodePackagesResponse struct {
	Packages []models.Package
}

type KillNodeProcessRequest struct {
	PID int32
}

type NodeSpeedtestResponse struct {
	Speedtest models.Speedtest
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
	gob.Register(RequiresAuthenticationResponse{})
	gob.Register(AuthenticateNetworkRequest{})
	gob.Register(AuthenticateNetworkResponse{})
	gob.Register(UpdateNodeStatsRequest{})
	gob.Register(UpdateNodeStatsPollDurationRequest{})
	gob.Register(NodeConnectionsResponse{})
	gob.Register(NodePackagesResponse{})
	gob.Register(NodeProcessesResponse{})
	gob.Register(KillNodeProcessRequest{})
	gob.Register(NodeSpeedtestResponse{})

	gob.Register(models.Node{})
	gob.Register(OPError{})

}
