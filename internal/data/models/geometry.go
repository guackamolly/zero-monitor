package models

import (
	"fmt"
	"time"
)

// Represents a point in a 2D space.
type Point struct {
	X float64
	Y float64
}

type Marker interface {
	X() float64
	Y() float64
	Float64() float64
	Label() string
}

// Defines a boundary in a 2D space.
type Axis struct {
	Min     Point
	Max     Point
	Markers []Marker
}

// Traces the points that compose a line in a 2D space.
type Line struct {
	Points []Marker
}

type ValueMarker struct {
	Point Point
	label string
	value float64
}

func (m ValueMarker) X() float64 {
	return m.Point.X
}

func (m ValueMarker) Y() float64 {
	return m.Point.Y
}

func (m ValueMarker) Float64() float64 {
	return m.value
}

func (m ValueMarker) Label() string {
	return m.label
}

func (a Axis) Fit(value float64) float64 {
	m1 := a.Markers[0]
	m2 := a.Markers[len(a.Markers)-1]

	// if X axis
	if a.Min.Y == a.Max.Y {
		return Lerp(m1.X(), m1.Float64(), m2.X(), m2.Float64(), value)
	}

	return Lerp(m1.Float64(), m1.Y(), m2.Float64(), m2.Y(), value)
}

func NewPoint(x, y float64) Point {
	return Point{
		X: x,
		Y: y,
	}
}

func NewAxis(
	min Point,
	max Point,
	labels int,
	gen func(int, float64, float64) Marker,
) Axis {
	if labels <= 2 {
		return Axis{
			Min: min,
			Max: max,
			Markers: []Marker{
				gen(0, min.X, min.Y),
				gen(1, max.X, max.Y),
			},
		}
	}

	markers := make([]Marker, labels)
	xdistance := (max.X - min.X) / float64(labels)
	ydistance := (max.Y - min.Y) / float64(labels)

	for i := range markers {
		if i == 0 {
			markers[i] = gen(
				i,
				0,
				0,
			)
			continue
		}

		markers[i] = gen(
			i,
			markers[i-1].X()+xdistance,
			markers[i-1].Y()+ydistance,
		)
	}

	for _, m := range markers {
		fmt.Printf("m.X(): %v\n", m.X())
	}

	return Axis{
		Min:     min,
		Max:     max,
		Markers: markers,
	}
}

func NewLine(
	points []Marker,
) Line {
	return Line{
		Points: points,
	}
}

func NewValueMarker(
	x, y float64,
	value float64,
	label string,
) ValueMarker {
	return ValueMarker{
		Point: NewPoint(x, y),
		label: label,
		value: value,
	}
}

func NewTimeAxis(
	min Point,
	max Point,
	values []time.Time,
) Axis {
	return NewAxis(
		min,
		max,
		2,
		func(i int, x, y float64) Marker {
			return NewValueMarker(
				x, y,
				float64(values[i].UnixMilli()),
				values[i].Format(time.TimeOnly),
			)
		},
	)
}

func Lerp(
	x1, y1 float64,
	x2, y2 float64,
	y float64,
) float64 {
	return x1 + (y-y1)*(x2-x1)/(y2-y1)
}
