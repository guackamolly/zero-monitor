package repositories

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SystemRepository interface {
	Info() (models.Info, error)
	Stats() (models.Stats, error)
}

type GopsUtilSystemRepository struct{}

func (r GopsUtilSystemRepository) Info() (models.Info, error) {
	cc, err := cpu.Counts(false)
	if err != nil {
		return models.Info{}, err
	}

	rs, err := mem.VirtualMemory()
	if err != nil {
		return models.Info{}, err
	}

	dsk, err := disk.Usage("/")
	if err != nil {
		return models.Info{}, err
	}

	os, err := host.Info()
	if err != nil {
		return models.Info{}, err
	}

	var cpuinfo cpu.InfoStat
	cpus, err := cpu.Info()
	if err != nil {
		logging.LogWarning("couldn't fetch cpu info, %v", err)
	}

	if len(cpus) > 0 {
		cpuinfo = cpus[0]
	}

	lip, err := localIP()
	if err != nil {
		logging.LogWarning("couldn't fetch local ip, %v", err)
		lip = []byte{}
	}

	pip, err := publicIP()
	if err != nil {
		logging.LogWarning("couldn't fetch public ip, %v", err)
		pip = []byte{}
	}

	return models.NewInfo(
		os.KernelArch,
		cpuinfo.ModelName,
		cpuinfo.CacheSize,
		cc,
		rs.Total,
		dsk.Total,
		os.Hostname,
		os.OS,
		fmt.Sprintf("%s %s", os.Platform, os.PlatformVersion),
		os.KernelVersion,
		lip,
		pip,
	), nil
}

func (r GopsUtilSystemRepository) Stats() (models.Stats, error) {
	cp, err := cpu.Percent(time.Millisecond*150, false)
	if err != nil {
		return models.Stats{}, err
	}

	rs, err := mem.VirtualMemory()
	if err != nil {
		return models.Stats{}, err
	}

	disk, err := disk.Usage("/")
	if err != nil {
		return models.Stats{}, err
	}

	st, _ := host.SensorsTemperatures()

	temp := float64(0)
	for _, t := range st {
		temp += t.Temperature
	}

	if l := len(st); l > 0 {
		temp = temp / float64(l)
	} else {
		temp = -1
	}

	uptime, err := host.Uptime()
	if err != nil {
		uptime = 0
	}

	return models.NewStats(
		cp[0],
		rs.UsedPercent,
		disk.UsedPercent,
		temp,
		uptime,
	), nil
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
