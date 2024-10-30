package repositories

import (
	"fmt"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	nett "github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type GopsUtilSystemRepository struct {
	TotalRx uint64
	TotalTx uint64
	Rx      uint64
	Tx      uint64
	procs   map[int32]models.Process
}

func NewGopsUtilSystemRepository() *GopsUtilSystemRepository {
	r := &GopsUtilSystemRepository{
		procs: map[int32]models.Process{},
	}

	go func() {
		for {
			time.Sleep(time.Second)
			ctr, err := nett.IOCounters(false)
			if err != nil {
				logging.LogError("failed to get IO counters, %v", err)
				continue
			}

			if len(ctr) == 1 {
				r.Rx = ctr[0].BytesRecv - r.TotalRx
				r.Tx = ctr[0].BytesSent - r.TotalTx
				r.TotalRx = ctr[0].BytesRecv
				r.TotalTx = ctr[0].BytesSent
			}
		}
	}()

	return r
}

func (r GopsUtilSystemRepository) Info() (models.MachineInfo, error) {
	cpucount, err := cpu.Counts(false)
	if err != nil {
		return models.MachineInfo{}, err
	}

	rs, err := mem.VirtualMemory()
	if err != nil {
		return models.MachineInfo{}, err
	}

	os, err := host.Info()
	if err != nil {
		return models.MachineInfo{}, err
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

	dsks := []models.DiskInfo{}
	block, err := ghw.Block()
	if err != nil {
		logging.LogWarning("couldn't fetch disks info, %v", err)
		block = &ghw.BlockInfo{}
	}
	for _, d := range block.Disks {
		info := models.NewDiskInfo(d.SizeBytes, d.Model, d.DriveType.String(), d.StorageController.String())
		dsks = append(dsks, info)
	}

	gpus := []models.GPUInfo{}
	gpu, err := ghw.GPU()
	if err != nil {
		logging.LogWarning("couldn't fetch gpus info, %v", err)
		gpu = &ghw.GPUInfo{}
	}
	for _, d := range gpu.GraphicsCards {
		info := models.NewGPUInfo(d.DeviceInfo.Product.Name, d.DeviceInfo.Vendor.Name)
		gpus = append(gpus, info)
	}

	product, err := ghw.Product()
	if err != nil {
		logging.LogWarning("couldn't fetch product info, %v", err)
		product = &ghw.ProductInfo{}
	}

	return models.NewMachineInfo(
		models.NewCPUInfo(
			os.KernelArch,
			cpuinfo.ModelName,
			cpuinfo.CacheSize,
			cpucount,
		),
		models.NewRAMInfo(
			rs.Total,
		),
		models.NewNetworkInfo(
			lip, pip,
		),
		models.NewOSInfo(
			os.Hostname,
			os.OS,
			fmt.Sprintf("%s %s", os.Platform, os.PlatformVersion),
			os.KernelVersion,
		),
		models.NewProductInfo(product.Family, product.Vendor),
		dsks,
		gpus,
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
		r.Rx,
		r.Tx,
		r.TotalRx,
		r.TotalTx,
	), nil
}

func (r GopsUtilSystemRepository) Conns() ([]models.Connection, error) {
	conns, err := nett.Connections("inet")
	if err != nil {
		return nil, err
	}

	l := len(conns)
	v := make([]models.Connection, l)
	for i, conn := range conns {
		v[i] = models.NewConnection(
			conn.Type,
			conn.Status,
			conn.Laddr,
			conn.Raddr,
		)
	}

	return v, nil
}

func (r GopsUtilSystemRepository) Procs() ([]models.Process, error) {
	pids, err := process.Pids()
	if err != nil {
		return nil, err
	}

	l := len(pids)
	v := make([]models.Process, l)
	cache := map[int32]models.Process{}

	for i, pid := range pids {
		pproc, err := process.NewProcess(pid)
		if err != nil {
			logging.LogError("failed to get data of process %d, %v", pid, err)
			continue
		}

		proc := r.cacheProc(pproc)
		mem, err := pproc.MemoryInfo()
		if err != nil {
			logging.LogWarning("failed to get memory info of process %d, %v", pproc.Pid, err)
			mem = &process.MemoryInfoStat{}
		}

		cpu, err := pproc.CPUPercent()
		if err != nil {
			logging.LogWarning("failed to get cpu info of process %d, %v", pproc.Pid, err)
			cpu = 0
		}

		cache[pid] = proc
		v[i] = proc.WithUpdatedMemory(mem.RSS).WithUpdatedCPU(cpu)
	}

	// delete processes that don't exist anymore
	for pid := range r.procs {
		if _, ok := cache[pid]; !ok {
			delete(r.procs, pid)
		}
	}

	return v, nil
}

func (r GopsUtilSystemRepository) KillProc(pid int32) error {
	_, ok := r.procs[pid]
	if !ok {
		return fmt.Errorf("process not associated to PID %d", pid)
	}

	pproc, err := process.NewProcess(pid)
	if err != nil {
		return err
	}

	err = pproc.Kill()
	if err != nil {
		return err
	}

	delete(r.procs, pid)
	return nil
}

func (r GopsUtilSystemRepository) cacheProc(pproc *process.Process) models.Process {
	pid := pproc.Pid
	if proc, ok := r.procs[pid]; ok {
		return proc
	}

	n, err := pproc.Name()
	if err != nil {
		logging.LogWarning("failed to get name of process %d, %v", pid, err)
		n = "-"
	}

	usr, err := pproc.Username()
	if err != nil {
		logging.LogWarning("failed to get user of process %d, %v", pid, err)
		usr = "-"
	}

	cmd, err := pproc.Cmdline()
	if err != nil {
		logging.LogWarning("failed to get command path of process %d, %v", pid, err)
		cmd = "-"
	}

	proc := models.NewProcess(
		pid,
		usr,
		n,
		cmd,
	)

	r.procs[pid] = proc
	return proc
}
