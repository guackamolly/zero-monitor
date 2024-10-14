package main

import (
	"log"

	"github.com/guackamolly/zero-monitor/internal/conn"
)

func main() {
	err := conn.StartBeaconServer()
	if err != nil {
		log.Fatalf("failed to start beacon server, %v", err)
	}
}
