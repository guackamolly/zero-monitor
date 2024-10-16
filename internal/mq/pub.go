package mq

import (
	"log"

	"github.com/guackamolly/zero-monitor/internal/di"
)

func (s Socket) RegisterPublishers() {
	pc := di.ExtractPublishContainer(s.ctx)
	if pc == nil {
		log.Fatalln("publish container hasn't been injected")
	}

	go func() {
		nr := pc.NodeReporter
		ns := nr.Start()
		n := <-ns

		s.PublishAndForget(compose(join, n))
		for n = range ns {
			s.PublishAndForget(compose(update, n))
		}
	}()
}
