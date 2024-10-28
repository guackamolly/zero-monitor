package mq

import (
	"context"

	"github.com/guackamolly/zero-monitor/internal/domain"
)

type containerKey int

const (
	keySubscribeContainer containerKey = iota
	keyPublishContainer
)

// Container for all dependencies required in a subscription context.
type SubscribeContainer struct {
	JoinNodesNetwork                   domain.JoinNodesNetwork
	UpdateNodesNetwork                 domain.UpdateNodesNetwork
	GetNodeStatsPollingDuration        domain.GetNodeStatsPollingDuration
	GetNodeStatsPollingDurationUpdates domain.GetNodeStatsPollingDurationUpdates
}

// Container for all dependencies required in a publish context.
type PublishContainer struct {
	GetCurrentNode            domain.GetCurrentNode
	GetCurrentNodeConnections domain.GetCurrentNodeConnections
	GetCurrentNodeProcesses   domain.GetCurrentNodeProcesses
	StartNodeStatsPolling     domain.StartNodeStatsPolling
	UpdateNodeStatsPolling    domain.UpdateNodeStatsPolling
	KillNodeProcess           domain.KillNodeProcess
}

func InjectSubscribeContainer(
	ctx context.Context,
	container SubscribeContainer,
) context.Context {
	return context.WithValue(ctx, keySubscribeContainer, &container)
}

func ExtractSubscribeContainer(
	ctx context.Context,
) *SubscribeContainer {
	return ctx.Value(keySubscribeContainer).(*SubscribeContainer)
}

func InjectPublishContainer(
	ctx context.Context,
	container PublishContainer,
) context.Context {
	return context.WithValue(ctx, keyPublishContainer, &container)
}

func ExtractPublishContainer(
	ctx context.Context,
) *PublishContainer {
	return ctx.Value(keyPublishContainer).(*PublishContainer)
}
