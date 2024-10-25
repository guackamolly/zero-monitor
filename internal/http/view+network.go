package http

import (
	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type NetworkNodeInformationView struct {
	NodeView
}

type NetworkNodeConnectionsView struct {
	NodeView
	Connections []models.Connection
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
