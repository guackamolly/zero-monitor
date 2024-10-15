package conn

import (
	"log"
	"net"
)

type Connection struct {
	IsMaster bool
	Port     int
	IP       net.IP
}

// starts the master node beacon server, which listens for
// UDP probe beacon requests on the local network.
func StartBeaconServer(conn *net.UDPConn) {
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

			d := decode(buf)
			if d.key != probeKey {
				continue
			}
			log.Printf("received probe beacon from %v\n", addr)

			go func() {
				b := encode(compose(helloKey))
				_, err = conn.WriteToUDP(b, addr)
				if err != nil {
					log.Printf("error replying to probe beacon, %v\n", err)
				}
			}()
		}
	}()
}
