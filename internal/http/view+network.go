package http

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type NetworkNodeInformationView struct {
	models.Node
}

type NetworkNodeConnectionsView struct {
	models.Node
	Connections []models.Connection
}

func NewNetworkNodeInformationView(
	node models.Node,
) NetworkNodeInformationView {
	return NetworkNodeInformationView{
		Node: node,
	}
}

func NewNetworkNodeConnectionsView(
	node models.Node,
	connections []models.Connection,
) NetworkNodeConnectionsView {
	return NetworkNodeConnectionsView{
		Node:        node,
		Connections: connections,
	}
}

func (v NetworkNodeInformationView) CPU() string {
	if len(v.Info.CPUModel) > 0 {
		return fmt.Sprintf("%s, %s, %d cores, %s cache", v.Info.CPUModel, v.Info.CPUArch, v.Info.CPUCount, v.Info.CPUCache)
	}

	return fmt.Sprintf("%s, %d cores", v.Info.CPUArch, v.Info.CPUCount)
}
