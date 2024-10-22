package http

import "fmt"

func (v NetworkNodeInformationView) CPU() string {
	if len(v.Info.CPUModel) > 0 {
		return fmt.Sprintf("%s, %s, %d cores, %s cache", v.Info.CPUModel, v.Info.CPUArch, v.Info.CPUCount, v.Info.CPUCache)
	}

	return fmt.Sprintf("%s, %d cores", v.Info.CPUArch, v.Info.CPUCount)
}
