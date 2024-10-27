package repositories

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type SystemRepository interface {
	Info() (models.MachineInfo, error)
	Stats() (models.Stats, error)
	Conns() ([]models.Connection, error)
	Procs() ([]models.Process, error)
}

func localIP() (net.IP, error) {
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

func publicIP() (net.IP, error) {
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
