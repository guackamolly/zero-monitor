package models_test

import (
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type lerpInput struct {
	x1, y1 float64
	x2, y2 float64
	y      float64
}

func TestLerp(t *testing.T) {

	testCases := []struct {
		desc   string
		input  lerpInput
		output float64
	}{
		{
			desc: "returns 1 if first point is (0,0), second point is (2, 2) and y=1",
			input: lerpInput{
				x1: 0,
				y1: 0,
				x2: 2,
				y2: 2,
				y:  1,
			},
			output: 1,
		},
		{
			desc: "returns 45 if first point is (40,0), second point is (50, 250) and y=125",
			input: lerpInput{
				x1: 40,
				y1: 0,
				x2: 50,
				y2: 250,
				y:  125,
			},
			output: 45,
		},
		{
			desc: "returns 125 if first point is (0,40), second point is (250, 50) and y=45",
			input: lerpInput{
				x1: 0,
				y1: 40,
				x2: 250,
				y2: 50,
				y:  45,
			},
			output: 125,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output := models.Lerp(tC.input.x1, tC.input.y1, tC.input.x2, tC.input.y2, tC.input.y)
			if output != tC.output {
				t.Errorf("expected %f but got %f", tC.output, output)
			}
		})
	}
}
