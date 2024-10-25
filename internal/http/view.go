package http

import "github.com/guackamolly/zero-monitor/internal/data/models"

type ServerStatsView struct {
	Online  []NodeView
	Offline []NodeView
}

type SettingsView struct {
	Form  FormView
	Error error
}

func NewServerStatsView(
	nodes []models.Node,
) ServerStatsView {
	on := []NodeView{}
	off := []NodeView{}
	for _, v := range nodes {
		if v.Online {
			on = append(on, NodeView(v))
		} else {
			off = append(off, NodeView(v))
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
