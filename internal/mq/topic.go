package mq

// Enumerates all publish and publish-reply topics.
type Topic byte

const (
	join Topic = iota + 1
	update
	reply
	empty
	xerror
	unknown
)
