package mq

import (
	"fmt"
	"os"
)

const (
	mqInviteCodeEnvKey = "mq_invite_code"
)

var (
	mqInviteCode = os.Getenv(mqInviteCodeEnvKey)
)

// This flag controls whether or not the current running node has already handshaked with the master node.
// It's used to disallow handshaking more than one time.
var handshaked = false

// Returns the invite code passed through environment variables ([mqInviteCodeEnvKey]) or reads it from std input.
func InviteCode() string {
	if len(mqInviteCode) == 0 {
		println("Waiting for invite code... (press enter to resume)")
		fmt.Scanln(&mqInviteCode)
	}

	return mqInviteCode
}
