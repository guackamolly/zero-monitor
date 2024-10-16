package models

import (
	"fmt"
)

const kb = 1024
const mb = kb * kb
const gb = mb * kb
const tb = gb * kb

type Percent float64
type Celsius float64
type Memory uint64

func (v Percent) String() string {
	return fmt.Sprintf("%0.2f%%", v)
}

func (v Celsius) String() string {
	return fmt.Sprintf("%f ºC", v)
}

func (v Memory) String() string {
	if v < kb {
		return fmt.Sprintf("%d B", v)
	}

	if v < mb {
		return fmt.Sprintf("%d KB", v/kb)
	}

	if v < gb {
		return fmt.Sprintf("%d MB", v/mb)
	}

	if v < tb {
		return fmt.Sprintf("%d GB", v/gb)
	}

	return fmt.Sprintf("%d TB", v/tb)
}
