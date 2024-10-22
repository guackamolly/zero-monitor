package http

import "github.com/guackamolly/zero-monitor/internal/data/models"

type ServerStatsView struct {
	Online  []models.Node
	Offline []models.Node
}

type NetworkNodeInformationView struct {
	models.Node
}

type SettingsView struct {
	Form  FormView
	Error error
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

func NewSettingsView(
	form FormView,
	error error,
) SettingsView {
	return SettingsView{
		Form:  form,
		Error: error,
	}
}

func NewNetworkNodeInformationView(
	node models.Node,
) NetworkNodeInformationView {
	return NetworkNodeInformationView{
		Node: node,
	}
}
