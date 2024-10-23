package models

import (
	"fmt"
	"net"
	"time"
)

const (
	kb = 1024
	mb = kb * kb
	gb = mb * kb
	tb = gb * kb
)

type Percent float64
type Celsius float64
type Memory uint64
type Duration time.Duration
type IP []byte

func (v Percent) String() string {
	if v < 0 {
		return "-"
	}

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

func (v Duration) String() string {
	return fmt.Sprint(time.Duration(v))
}

func (v IP) String() string {
	if len(v) < 4 {
		return "-"
	}

	return net.IP(v).String()
}
