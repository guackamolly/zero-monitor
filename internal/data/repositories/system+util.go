package repositories

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
)

// Returns the private IP of the first interface found that is not a loopback interface.
//
// Err is not nil if:
// 1) interfaces lookup call fails
// 2) all interfaces are loopbacks
func InterfaceIP() (net.IP, error) {
	iaddrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, iaddr := range iaddrs {
		c, ok := iaddr.(*net.IPNet)
		if !ok {
			continue
		}

		if c.IP.IsLoopback() || !c.IP.IsPrivate() {
			continue
		}

		return c.IP, nil
	}

	return nil, fmt.Errorf("couldn't lookup interface IP")
}

// Returns the private IP of the preferred interface for remote network requests
// by submitting an UDP packet to 1.1.1.1.
//
// Err is not nil if UDP packet submission fails.
func PrivateIP() (net.IP, error) {
	conn, err := net.Dial("udp4", "1.1.1.1:80")
	if err != nil {
		return nil, err
	}

	laddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return nil, fmt.Errorf("failed to cast address")
	}

	return laddr.IP, conn.Close()
}

// Returns the host public IP (IPv4) by using an external lookup service (GET ifconfig.me).
//
// Err is not nil if the HTTP GET call fails.
func PublicIP() (net.IP, error) {
	myDialer := net.Dialer{}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return myDialer.DialContext(ctx, "tcp4", addr)
	}

	client := http.Client{
		Transport: transport,
	}
	resp, err := client.Get("http://ifconfig.me")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(bs))
	if ip == nil {
		return nil, fmt.Errorf("failed to parse %s to IP", bs)
	}

	ip4 := ip.To4()
	if ip4 != nil {
		return ip4, nil
	}

	return ip, nil
}
