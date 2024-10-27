package mq

// Enumerates all publish and publish-reply topics.
type Topic byte

const (
	JoinNetwork Topic = iota + 1
	UpdateNodeStats
	UpdateNodeStatsPollDuration
	NodeConnections
	NodeProcesses
)
