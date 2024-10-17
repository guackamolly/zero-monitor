package http

import "github.com/guackamolly/zero-monitor/internal/data/models"

type ServerStatsView struct {
	Online  []models.Node
	Offline []models.Node
}

func NewServerStatsView(
	nodes []models.Node,
) ServerStatsView {
	on := []models.Node{}
	off := []models.Node{}
	for _, v := range nodes {
		if v.Online {
			on = append(on, v)
		} else {
			off = append(off, v)
		}
	}

	return ServerStatsView{
		Online:  on,
		Offline: off,
	}
}
