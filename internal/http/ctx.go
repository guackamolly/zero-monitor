package http

import (
	"context"

	"github.com/guackamolly/zero-monitor/internal/service"
)

type containerKey int

const (
	keyServiceContainer containerKey = iota
)

// Container for all dependencies required in service context.
type ServiceContainer struct {
	NodeManager         *service.NodeManagerService
	NodeScheduler       *service.NodeSchedulerService
	NodeCommander       *service.NodeCommanderService
	NodeSpeedtest       *service.NodeSpeedtestService
	MasterConfiguration *service.MasterConfigurationService
	Network             *service.NetworkService
	Networking          *service.NetworkingService
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
