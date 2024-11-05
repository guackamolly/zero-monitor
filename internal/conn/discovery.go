package conn

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

// At most a node can send 5 beacon requests in parallel.
const maxParallelBeacons = 5

// A node must wait at most 2 seconds for a response from the master node.
// Two seconds is more than enough for a reply since beacons are sent in broadcast on the local network.
const beaconReplyWaitDuration = 2 * time.Second

var NetworkIP = net.IPv4(0, 0, 0, 0)
var broadcastIP = net.IPv4(255, 255, 255, 255)

// starts broadcasting beacon probes in the local network, on all known ports
// that a master node could be registered in.
func StartBeaconBroadcast() (Connection, error) {
	type result struct {
		conn Connection
		err  error
	}

	c := make(chan (result), 1)
	stop := false
	go func() {
		pb := 0
		for _, p := range ports {
			pb++

			if pb > maxParallelBeacons {
				time.Sleep(beaconReplyWaitDuration)
				pb = pb - maxParallelBeacons
			}

			if stop {
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
	close(c)
	stop = true

	return r.conn, r.err
}

// Finds a port that is available for incoming TCP connections.
// If err is nil, then an open connection is returned, as a way
// to reserve that same TCP port. This connection must be closed
// before the port is used in a different context.
func FindAvailableTcpPort(
	ip net.IP,
) (*net.TCPListener, error) {
	var conn *net.TCPListener
	var err error

	// first find a port from the list of well known ports, that is available for listening
	for _, p := range ports {
		addr := net.TCPAddr{Port: int(p), IP: ip}
		conn, err = net.ListenTCP("tcp", &addr)

		if err == nil {
			break
		}
	}

	// if no port from the list of well known ports is available for listening, connect to a random one.
	if err != nil {
		addr := net.TCPAddr{Port: 0, IP: ip}
		conn, err = net.ListenTCP("tcp", &addr)
	}

	return conn, err
}

// Same thing as [FindAvailableTcpPort], but for UDP connections.
func FindAvailableUdpPort(
	ip net.IP,
) (*net.UDPConn, error) {
	var conn *net.UDPConn
	var err error

	// first find a port from the list of well known ports, that is available for listening
	for _, p := range ports {
		addr := net.UDPAddr{Port: int(p), IP: ip}
		conn, err = net.ListenUDP("udp", &addr)

		if err == nil {
			break
		}
	}

	// if no port from the list of well known ports is available for listening, connect to a random one.
	if err != nil {
		addr := net.UDPAddr{Port: 0, IP: ip}
		conn, err = net.ListenUDP("udp", &addr)
	}

	return conn, err
}

// sends a broadcast probe beacon and waits for a response
func broadcastProbeBeacon(port uint16) (Connection, error) {
	addr := net.UDPAddr{IP: broadcastIP, Port: int(port)}
	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		return Connection{}, err
	}

	defer conn.Close()
	bs, err := models.Encode(compose(probeKey))
	if err != nil {
		return Connection{}, err
	}

	_, err = conn.Write(bs)
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
	defer conn.Close()

	// set 10 sec timeout until the connection is closed
	// since these beacon messages are only sent on the local network,
	// then 10 seconds is more than enough more the master node to answer
	err = conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return Connection{}, err
	}

	buf := make([]byte, 1024)
	_, _, err = conn.ReadFromUDP(buf)
	if err != nil {
		return Connection{}, err
	}

	d, err := models.Decode[msg](buf)
	if err != nil {
		return Connection{}, err
	}

	if d.Key != helloKey {
		return Connection{}, fmt.Errorf("received unknown response after sending probe beacon, %v", d)
	}

	subConn, ok := d.Data.(Connection)
	if !ok {
		return Connection{}, fmt.Errorf("couldn't parse data to Connection struct, %v", d.Data)
	}

	return subConn, nil
}
