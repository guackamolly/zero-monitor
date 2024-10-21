package mq

import (
	"fmt"
	"log"

	"github.com/guackamolly/zero-monitor/internal/di"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/guackamolly/zero-monitor/internal/service"
)

func (s Socket) RegisterPublishers() {
	pc := di.ExtractPublishContainer(s.ctx)
	if pc == nil {
		log.Fatalln("publish container hasn't been injected")
	}

	go func() {
		// 1. Join network.
		nr := pc.NodeReporter

		resp, err := joinNetwork(s, nr)
		if err != nil {
			// TODO: handle join network error gracefully.
			logging.LogFatal("couldn't join network, %v", err)
		}

		// 2. Start stats polling.
		ns := nr.Start(resp.StatsPoll)
		for n := range ns {
			s.PublishAndForget(compose(update, n))
		}
	}()
}

func joinNetwork(
	s Socket,
	nodeReporter *service.NodeReporterService,
) (joinNodeResponse, error) {
	n := nodeReporter.Initial()
	m, err := s.Publish(compose(join, joinNodeRequest{Node: n}))
	if err != nil {
		return joinNodeResponse{}, err
	}

	if m.Topic == xerror {
		return joinNodeResponse{}, fmt.Errorf("%v", m.Data)
	}

	if m.Topic != reply {
		return joinNodeResponse{}, fmt.Errorf("couldn't understand reply message topic: %v, data: %v", m.Topic, m.Data)
	}

	resp, ok := m.Data.(joinNodeResponse)
	if !ok {
		return joinNodeResponse{}, fmt.Errorf("couldn't parse data to join node response, %v", m.Data)
	}

	return resp, nil
}
