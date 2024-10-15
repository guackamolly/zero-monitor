package models

type Node struct {
	ID    string
	Info  Info
	Stats Stats
}

func (m Node) WithUpdatedStats(stats Stats) Node {
	return Node{
		ID:    m.ID,
		Info:  m.Info,
		Stats: stats,
	}
}
