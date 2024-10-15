package models

type Node struct {
	ID    string
	Info  Info
	Stats Stats
}

func NewNode(
	id string,
	info Info,
	stats Stats,
) Node {
	return Node{
		ID:    id,
		Info:  info,
		Stats: stats,
	}
}

func NewNodeWithoutStats(
	id string,
	info Info,
) Node {
	stats := UnknownStats()
	return NewNode(id, info, stats)
}

func (m Node) WithUpdatedStats(stats Stats) Node {
	return Node{
		ID:    m.ID,
		Info:  m.Info,
		Stats: stats,
	}
}
