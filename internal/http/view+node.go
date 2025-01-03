package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type NodeView models.Node

func (v NodeView) Hostname() string {
	return v.Info.OS.Hostname
}

func (v NodeView) OS() string {
	return v.Info.OS.OS
}

func (v NodeView) OSType() string {
	os := v.OS()
	distro := v.Distribution()
	if os != "linux" {
		return os
	}

	if strings.HasPrefix(distro, "raspbian") {
		return "raspbian"
	}

	return os
}

func (v NodeView) Distribution() string {
	return v.Info.OS.Distribution
}

func (v NodeView) Kernel() string {
	return v.Info.OS.Kernel
}

func (v NodeView) TotalRAM() string {
	return v.Info.RAM.Total.String()
}

func (v NodeView) TotalDisk() string {
	if len(v.Info.Disk) == 0 {
		return "-"
	}

	return v.Info.Disk[0].Total.String()
}

func (v NodeView) Uptime() string {
	return v.Stats.Uptime.String()
}

func (v NodeView) CPUCount() string {
	return fmt.Sprintf("%d", v.Info.CPU.Count)
}

func (v NodeView) CPUUsage() string {
	return v.Stats.CPU.String()
}

func (v NodeView) CPUUsageLevel() string {
	return fillLevel(v.Stats.CPU)
}

func (v NodeView) RAMUsage() string {
	return v.Stats.RAM.String()
}

func (v NodeView) RAMUsageLevel() string {
	return fillLevel(v.Stats.RAM)
}

func (v NodeView) DiskUsage() string {
	return v.Stats.Disk.String()
}

func (v NodeView) DiskUsageLevel() string {
	return fillLevel(v.Stats.Disk)
}

func (v NodeView) LocalIP() string {
	return v.Info.Network.LocalIP.String()
}

func (v NodeView) PublicIP() string {
	return v.Info.Network.PublicIP.String()
}

func (v NodeView) IsSingleDisk() bool {
	return len(v.Info.Disk) == 1
}

func (v NodeView) HasDisk() bool {
	return len(v.Info.Disk) > 0
}

func (v NodeView) HasGPU() bool {
	return len(v.Info.GPU) > 0
}

func (v NodeView) DiskCount() int {
	return len(v.Info.Disk)
}

func (v NodeView) GPUCount() int {
	return len(v.Info.GPU)
}

func (v NodeView) CPU() string {
	cpu := v.Info.CPU
	if len(cpu.Model) > 0 {
		return fmt.Sprintf("%s, %s, %d cores, %s cache", cpu.Model, cpu.Arch, cpu.Count, cpu.Cache)
	}

	return fmt.Sprintf("%s, %d cores", cpu.Arch, cpu.Count)
}

func (v NodeView) RAM() string {
	return v.TotalRAM()
}

func (v NodeView) Disk(idx int) string {
	if len(v.Info.Disk) == 0 {
		return "-"
	}

	dsk := v.Info.Disk[idx]
	return fmt.Sprintf("%s %s - %s (%s %s)", dsk.Type, dsk.Kind, dsk.Model, dsk.Name, dsk.Total)
}

func (v NodeView) GPU(idx int) string {
	if !v.HasGPU() {
		return "-"
	}

	gpu := v.Info.GPU[idx]
	return fmt.Sprintf("%s - %s", gpu.Vendor, gpu.Model)
}

func (v NodeView) Product() string {
	product := v.Info.Product
	if len(product.Model) == 0 && len(product.Vendor) == 0 {
		return "-"
	}

	return fmt.Sprintf("%s - %s", product.Vendor, product.Model)
}

func (v NodeView) LastSeenOn() string {
	return v.LastSeen.Format(time.DateTime)
}

func (v NodeView) Rx() string {
	return v.Stats.Rx.String()
}

func (v NodeView) Tx() string {
	return v.Stats.Tx.String()
}

func (v NodeView) TotalRx() string {
	return v.Stats.TotalRx.String()
}

func (v NodeView) TotalTx() string {
	return v.Stats.TotalTx.String()
}

func fillLevel(p models.Percent) string {
	if p < 0 {
		return "-"
	}

	if p < 60 {
		return "1"
	}

	if p < 85 {
		return "2"
	}

	return "3"
}
