package models

type MachineInfo struct {
	CPU     CPUInfo
	RAM     RAMInfo
	Network NetworkInfo
	OS      OSInfo
	Product ProductInfo
	Disk    []DiskInfo
	GPU     []GPUInfo
}

type CPUInfo struct {
	Arch  string
	Count int
	Model string
	Cache Memory
}

type RAMInfo struct {
	Total Memory
}

type DiskInfo struct {
	Total Memory
	Model string
	Type  string
	Kind  string
}

type GPUInfo struct {
	Model  string
	Vendor string
}

type OSInfo struct {
	Hostname     string
	OS           string
	Distribution string
	Kernel       string
}

type ProductInfo struct {
	Vendor string
	Model  string
}

type NetworkInfo struct {
	LocalIP  IP
	PublicIP IP
}

func NewCPUInfo(
	arch string,
	model string,
	cache int32,
	count int,
) CPUInfo {
	return CPUInfo{
		Arch:  arch,
		Count: count,
		Model: model,
		Cache: Memory(cache),
	}
}

func NewRAMInfo(
	total uint64,
) RAMInfo {
	return RAMInfo{
		Total: Memory(total),
	}
}

func NewDiskInfo(
	total uint64,
	model string,
	ttype string,
	kind string,
) DiskInfo {
	return DiskInfo{
		Total: Memory(total),
		Model: model,
		Type:  ttype,
		Kind:  kind,
	}
}

func NewGPUInfo(
	model string,
	vendor string,
) GPUInfo {
	return GPUInfo{
		Model:  model,
		Vendor: vendor,
	}
}

func NewOSInfo(
	hostname string,
	os string,
	distribution string,
	kernel string,
) OSInfo {
	return OSInfo{
		Hostname:     hostname,
		OS:           os,
		Distribution: distribution,
		Kernel:       kernel,
	}
}

func NewNetworkInfo(
	localIP []byte,
	publicIP []byte,
) NetworkInfo {
	return NetworkInfo{
		LocalIP:  localIP,
		PublicIP: publicIP,
	}
}

func NewProductInfo(
	model string,
	vendor string,
) ProductInfo {
	return ProductInfo{
		Model:  model,
		Vendor: vendor,
	}
}

func NewMachineInfo(
	cpu CPUInfo,
	ram RAMInfo,
	network NetworkInfo,
	os OSInfo,
	product ProductInfo,
	disks []DiskInfo,
	gpus []GPUInfo,
) MachineInfo {
	return MachineInfo{
		CPU:     cpu,
		RAM:     ram,
		Network: network,
		OS:      os,
		Product: product,
		Disk:    disks,
		GPU:     gpus,
	}
}
