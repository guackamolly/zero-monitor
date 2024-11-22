package models

import "time"

type JoinNetwork struct {
	NodeID   string
	NodeName string
	Code     string
}

type JoinNetworkCode struct {
	Code      string
	ExpiresAt time.Time
}

type NetworkNode struct {
	ID   string
	Name string
}

func (c JoinNetworkCode) Expired() bool {
	return c.ExpiresAt.Before(time.Now())
}

func NewJoinNetwork(
	id, name, code string,
) JoinNetwork {
	return JoinNetwork{
		NodeID:   id,
		NodeName: name,
		Code:     code,
	}
}

func NewJoinNetworkCode() JoinNetworkCode {
	return JoinNetworkCode{
		Code:      UUID(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

func NewNetworkNode(
	id, name string,
) NetworkNode {
	return NetworkNode{
		ID:   id,
		Name: name,
	}
}
