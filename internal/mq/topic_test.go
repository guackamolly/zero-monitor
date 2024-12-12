package mq_test

import (
	"testing"

	"github.com/guackamolly/zero-monitor/internal/mq"
)

func TestTopicSensitive(t *testing.T) {
	testCases := []struct {
		desc   string
		input  []mq.Topic
		output bool
	}{
		{
			desc:   "hello, goodbye, join, authenticate network, node connections, node processes, kill node process are sensitive",
			input:  []mq.Topic{mq.HelloNetwork, mq.GoodbyeNetwork, mq.JoinNetwork, mq.AuthenticateNetwork, mq.NodeConnections, mq.NodeProcesses, mq.KillNodeProcess},
			output: true,
		},
		{
			desc:   "everything else is not sensitive",
			input:  []mq.Topic{mq.UpdateNodeStats, mq.UpdateNodeStatsPollDuration, mq.StartNodeSpeedtest, mq.NodePackages},
			output: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			for _, tp := range tC.input {
				if output := tp.Sensitive(); output != tC.output {
					t.Errorf("expected %v but got %v", tC.output, output)
				}
			}
		})
	}
}
