package mq

import (
	"log"
)

func (s Socket) RegisterPublishers() {
	pc := ExtractPublishContainer(s.ctx)
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
