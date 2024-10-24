package di

import (
	"context"

	"github.com/guackamolly/zero-monitor/internal/domain"
	"github.com/guackamolly/zero-monitor/internal/service"
)

type containerKey int

const (
	keySubscribeContainer containerKey = iota
	keyPublishContainer
	keyServiceContainer
)

// Container for all dependencies required in service context.
type ServiceContainer struct {
	NodeManager         *service.NodeManagerService
	NodeScheduler       *service.NodeSchedulerService
	NodeCommander       *service.NodeCommanderService
	MasterConfiguration *service.MasterConfigurationService
}

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
	StartNodeStatsPolling     domain.StartNodeStatsPolling
	UpdateNodeStatsPolling    domain.UpdateNodeStatsPolling
}

func InjectServiceContainer(
	ctx context.Context,
	container ServiceContainer,
) context.Context {
	return context.WithValue(ctx, keyServiceContainer, &container)
}

func ExtractServiceContainer(
	ctx context.Context,
) *ServiceContainer {
	return ctx.Value(keyServiceContainer).(*ServiceContainer)
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
