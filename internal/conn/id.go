package conn

import (
	"log"

	"github.com/denisbrodbeck/machineid"
)

// Holds an unique machine identifier, hashed for security
var machineId string

func init() {
	mid, err := machineid.ProtectedID("zero-monitor")
	if err != nil {
		log.Fatalf("failed to extract machine id (required to identify nodes), %v\n", err)
	}

	machineId = mid
}
