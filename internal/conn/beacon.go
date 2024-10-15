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

var networkIP = net.IPv4(0, 0, 0, 0)

// starts the master node beacon server, which listens for
// UDP probe beacon requests on the local network.
func StartBeaconServer() error {
	var conn *net.UDPConn
	var err error

	// first find a port from the list of well known ports, that is available for listening
	for _, p := range ports {
		addr := net.UDPAddr{Port: int(p), IP: networkIP}
		conn, err = net.ListenUDP("udp", &addr)

		if err == nil {
			break
		}
	}

	if err != nil {
		return err
	}

	log.Printf("started udp beacon server on %v\n", conn.LocalAddr())
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
}
