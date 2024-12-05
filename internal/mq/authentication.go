package mq

// This flag controls whether or not the current running node has already handshaked with the master node.
// It's used to disallow handshaking more than one time.
var handshaked = false
