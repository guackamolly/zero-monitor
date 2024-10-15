package conn

import (
	"fmt"
	"log"
	"net"
	"time"
)

// At most a node can send 5 beacon requests in parallel.
const maxParallelBeacons = 5

// A node must wait at most 2 seconds for a response from the master node.
// Two seconds is more than enough for a reply since beacons are sent in broadcast on the local network.
const beaconReplyWaitDuration = 2 * time.Second

var broadcastIP = net.IPv4(255, 255, 255, 255)

// starts broadcasting beacon probes in the local network, on all known ports
// that a master node could be registered in.
func StartBeaconBroadcast() (Connection, error) {
	type result struct {
		conn Connection
		err  error
	}

	c := make(chan (result), 1)
	go func() {
		pb := 0
		for _, p := range ports {
			pb++

			if pb > maxParallelBeacons {
				time.Sleep(beaconReplyWaitDuration)
				pb = pb - maxParallelBeacons
			}

			if len(c) != 0 {
				return
			}

			go func() {
				log.Printf("sending broadcast beacon on port %d\n", p)
				conn, err := broadcastProbeBeacon(p)

				if err != nil {
					return
				}

				c <- result{conn: conn}
			}()
		}

		c <- result{err: fmt.Errorf("no master node found in all %d ports probed", len(ports))}
	}()

	r := <-c
	return r.conn, r.err
}

// sends a broadcast probe beacon and waits for a response
func broadcastProbeBeacon(port uint16) (Connection, error) {
	addr := net.UDPAddr{IP: broadcastIP, Port: int(port)}
	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		return Connection{}, err
	}

	defer conn.Close()
	_, err = conn.Write(encode(compose(probeKey)))
	if err != nil {
		return Connection{}, err
	}

	// close to local address
	conn.Close()

	laddr := conn.LocalAddr()
	conn, err = net.ListenUDP("udp", laddr.(*net.UDPAddr))
	if err != nil {
		return Connection{}, err
	}

	// set 2 sec timeout until the connection is closed
	// since these beacon messages are only sent on the local network,
	// then 2 seconds is more than enough more the master node to answer
	err = conn.SetDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return Connection{}, err
	}

	buf := make([]byte, 1024)
	var raddr *net.UDPAddr
	_, raddr, err = conn.ReadFromUDP(buf)
	if err != nil {
		return Connection{}, err
	}

	if d := decode(buf); d.key != helloKey {
		return Connection{}, fmt.Errorf("received unknown response after sending probe beacon, %v", d)
	}

	return Connection{
		Port:     raddr.Port,
		IP:       raddr.IP,
		IsMaster: true,
	}, nil
}
