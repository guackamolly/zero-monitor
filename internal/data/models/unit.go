package models

import (
	"fmt"
	"time"
)

const (
	kb = 1024
	mb = kb * kb
	gb = mb * kb
	tb = gb * kb

	Kbit = 1000.0
	Mbit = Kbit * Kbit
	Gbit = Mbit * Kbit

	m  = 1
	km = 1000 * m
)

type Percent float64
type Celsius float64
type Memory uint64
type Duration time.Duration
type IORate uint64
type BitRate float64
type Distance float64

func (v Percent) String() string {
	if v < 0 {
		return "-"
	}

	return fmt.Sprintf("%0.2f%%", v)
}

func (v Celsius) String() string {
	return fmt.Sprintf("%0.0f ºC", v)
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
	d := time.Duration(v)

	if d < time.Microsecond {
		return fmt.Sprintf("%d ns", v)
	}

	if d < time.Millisecond {
		return fmt.Sprintf("%d us", d.Microseconds())
	}

	if d < time.Second {
		return fmt.Sprintf("%d ms", d.Milliseconds())
	}

	if d < time.Minute {
		return fmt.Sprintf("%0.0f s", d.Seconds())
	}

	if d < time.Hour {
		return fmt.Sprintf("%0.0f min", d.Minutes())
	}

	return d.String()
}

func (v IORate) String() string {
	return fmt.Sprintf("%s/s", Memory(v))
}

func (v BitRate) String() string {
	if v < Kbit {
		return fmt.Sprintf("%0.1f bps", v)
	}

	if v < Mbit {
		return fmt.Sprintf("%0.1f Kbps", v/Kbit)
	}

	if v < Gbit {
		return fmt.Sprintf("%0.1f Mbps", v/Mbit)
	}

	return fmt.Sprintf("%0.1f Gbps", v/Gbit)
}

func (v Distance) String() string {
	if v < km {
		return fmt.Sprintf("%0.1f m", v)
	}

	return fmt.Sprintf("%0.1f Km", v/km)
}
