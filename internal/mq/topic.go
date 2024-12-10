package mq

// Enumerates all publish and publish-reply topics.
type Topic byte

const (
	HelloNetwork Topic = iota + 1
	JoinNetwork
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
	case HelloNetwork, JoinNetwork, AuthenticateNetwork, NodeConnections, NodeProcesses, KillNodeProcess:
		return true
	default:
		return false
	}
}
