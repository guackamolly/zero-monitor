package conn

import (
	"log"
	"net"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type Connection struct {
	Port int
	IP   net.IP
}

// starts the master node beacon server, which listens for
// UDP probe beacon requests on the local network.
func StartBeaconServer(conn *net.UDPConn, subConn Connection) {
	go func() {
		defer conn.Close()
		for {
			log.Println("waiting for probe requests...")
			buf := make([]byte, 1024)
			_, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("error reading UDP data, %v\n", err)
				continue
			}

			d, err := models.Decode[msg](buf)
			if err != nil {
				log.Printf("failed to decode beacon data, %v\n", err)
				continue
			}

			if d.Key != probeKey {
				continue
			}
			log.Printf("received probe beacon from %v\n", addr)

			go func() {
				b, err := models.Encode(compose(helloKey, subConn))
				if err != nil {
					log.Printf("failed to encode beacon reply data, %v", err)
					return
				}

				_, err = conn.WriteToUDP(b, addr)
				if err != nil {
					log.Printf("error replying to probe beacon, %v\n", err)
					return
				}
			}()
		}
	}()
}
