package di

import (
	"context"

	"github.com/guackamolly/zero-monitor/internal/service"
)

type containerKey int

const (
	keySubscribeContainer containerKey = iota
	keyPublishContainer
)

// Container for all dependencies required in a subscription context.
type SubscribeContainer struct {
	NodeManager         *service.NodeManagerService
	NodeScheduler       *service.NodeSchedulerService
	MasterConfiguration *service.MasterConfigurationService
}

// Container for all dependencies required in a publish context.
type PublishContainer struct {
	NodeReporter *service.NodeReporterService
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
