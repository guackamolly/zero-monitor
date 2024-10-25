package http

import (
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type NetworkView struct {
	Online             []NodeView
	Offline            []NodeView
	NodeStatsPollBurst time.Duration
}

type NetworkNodeInformationView struct {
	NodeView
}

type NetworkNodeConnectionsView struct {
	NodeView
	Connections []models.Connection
}

func NewNetworkView(
	nodes []models.Node,
	nodeStatsPollBurst time.Duration,
) NetworkView {
	on := []NodeView{}
	off := []NodeView{}
	for _, v := range nodes {
		if v.Online {
			on = append(on, NodeView(v))
		} else {
			off = append(off, NodeView(v))
		}
	}

	return NetworkView{
		Online:             on,
		Offline:            off,
		NodeStatsPollBurst: nodeStatsPollBurst,
	}
}

func NewNetworkNodeInformationView(
	node models.Node,
) NetworkNodeInformationView {
	return NetworkNodeInformationView{
		NodeView: NodeView(node),
	}
}

func NewNetworkNodeConnectionsView(
	node models.Node,
	connections []models.Connection,
) NetworkNodeConnectionsView {
	return NetworkNodeConnectionsView{
		NodeView:    NodeView(node),
		Connections: connections,
	}
}
