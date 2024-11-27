package mq

// Enumerates all publish and publish-reply topics.
type Topic byte

const (
	JoinNetwork Topic = iota + 1
	AuthenticateNetwork
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
	case JoinNetwork, AuthenticateNetwork, NodeConnections, NodeProcesses, KillNodeProcess:
		return true
	default:
		return false
	}
}
