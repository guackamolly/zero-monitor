package main

import (
	"log"

	"github.com/guackamolly/zero-monitor/internal/conn"
)

func main() {
	conn, err := conn.StartBeaconBroadcast()
	if err != nil {
		log.Fatalf("failed to probe master node, %v\n", err)
	}

	log.Printf("found master node, %v\n", conn)
}
