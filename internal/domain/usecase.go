package domain

import (
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type GetCurrentNode func() models.Node
type GetCurrentNodeConnections func() ([]models.Connection, error)
type GetCurrentNodeProcesses func() ([]models.Process, error)

type StartNodeStatsPolling func(d time.Duration) chan (models.Node)
type UpdateNodeStatsPolling func(d time.Duration)
type GetNodeStatsPollingDuration func() time.Duration
type GetNodeStatsPollingDurationUpdates func() chan (time.Duration)

type JoinNodesNetwork func(models.Node) error
type UpdateNodesNetwork func(models.Node) error

type KillNodeProcess func(int32) error
