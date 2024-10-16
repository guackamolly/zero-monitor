package http

import "github.com/guackamolly/zero-monitor/internal/data/models"

type ServerStatsView struct {
	Nodes []models.Node
}

func NewServerStatsView(
	nodes []models.Node,
) ServerStatsView {
	return ServerStatsView{
		Nodes: nodes,
	}
}
