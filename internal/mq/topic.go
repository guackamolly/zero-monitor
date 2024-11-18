package mq

// Enumerates all publish and publish-reply topics.
type Topic byte

const (
	JoinNetwork Topic = iota + 1
	UpdateNodeStats
	UpdateNodeStatsPollDuration
	NodeConnections
	NodeProcesses
	NodePackages
	KillNodeProcess
	StartNodeSpeedtest
)

func (t Topic) Sensitive() bool {
	switch t {
	case JoinNetwork, NodeConnections, NodeProcesses, KillNodeProcess:
		return true
	default:
		return false
	}
}
