package internal

import (
	"log"

	"github.com/denisbrodbeck/machineid"
)

// Holds an unique machine identifier, hashed for security
var MachineId string

func init() {
	mid, err := machineid.ProtectedID("zero-monitor")
	if err != nil {
		log.Fatalf("failed to extract machine id (required to identify nodes), %v\n", err)
	}

	MachineId = mid
}
