package models

import "time"

type Node struct {
	ID       string
	Online   bool
	LastSeen time.Time
	Info     MachineInfo
	Stats    Stats
}

func NewNode(
	id string,
	info MachineInfo,
	stats Stats,
) Node {
	return Node{
		ID:       id,
		Info:     info,
		Stats:    stats,
		Online:   true,
		LastSeen: time.Now(),
	}
}

func NewNodeWithoutStats(
	id string,
	info MachineInfo,
) Node {
	stats := UnknownStats()
	return NewNode(id, info, stats)
}

func (m Node) WithUpdatedStats(stats Stats) Node {
	return NewNode(m.ID, m.Info, stats)
}

func (m Node) SetOffline() Node {
	return Node{
		ID:       m.ID,
		Info:     m.Info,
		Stats:    m.Stats,
		LastSeen: m.LastSeen,
		Online:   false,
	}
}
