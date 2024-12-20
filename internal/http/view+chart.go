package http

import (
	"bytes"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/wcharczuk/go-chart/v2"
)

type ChartView interface {
	SVG(width, height int) string
}

type AxisView struct {
	Min    float64
	Max    float64
	Legend string
}

type LineView struct {
	XValues         []float64
	YValues         []float64
	XFormatter      func(float64) string
	YFormatter      func(float64) string
	TooltipProvider func(int) string
	Legend          string
}

type LineChartView struct {
	Lines []LineView
	X     AxisView
	Y     AxisView
}

func (v LineView) build() chart.ContinuousSeries {
	return chart.ContinuousSeries{
		XValues: v.XValues,
		XValueFormatter: func(i interface{}) string {
			return v.XFormatter(i.(float64))
		},
		YValues: v.YValues,
		YValueFormatter: func(i interface{}) string {
			return v.YFormatter(i.(float64))
		},
		Style: chart.Style{
			DotWidth:           4.0,
			DotTooltipProvider: v.TooltipProvider,
		},
		Name: v.Legend,
	}
}

func (v LineChartView) build(width, height int) chart.Chart {
	s := make([]chart.Series, len(v.Lines))
	for i, l := range v.Lines {
		s[i] = l.build()
	}

	cht := chart.Chart{
		Background: chart.Style{
			FillColor: chart.ColorTransparent,
		},
		Canvas: chart.Style{
			FillColor: chart.ColorTransparent,
		},

		Series: s,
		XAxis: chart.XAxis{
			Range: &chart.ContinuousRange{
				Min: v.X.Min,
				Max: v.X.Max,
			},
			Name: v.X.Legend,
		},
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Min: v.Y.Min,
				Max: v.Y.Max,
			},
			Name: v.Y.Legend,
		},
		YAxisSecondary: chart.HideYAxis(),
		Width:          width,
		Height:         height,
	}

	cht.Elements = []chart.Renderable{
		chart.Legend(&cht, chart.Style{
			FillColor: chart.ColorTransparent,
		}),
	}

	return cht
}

func (v LineChartView) SVG(width, height int) string {
	buffer := bytes.NewBuffer([]byte{})
	err := v.build(width, height).Render(chart.SVG, buffer)
	if err != nil {
		logging.LogWarning("failed to render chart as svg, %v", err)
	}

	return buffer.String()
}

func NewAxisView(
	min, max float64,
	legend string,
) AxisView {
	return AxisView{
		Min:    min,
		Max:    max,
		Legend: legend,
	}
}

func NewLineView(
	legend string,
	xvalues []float64,
	yvalues []float64,
	xformatter func(float64) string,
	yFormatter func(float64) string,
	tooltipProvider func(int) string,
) LineView {
	return LineView{
		XValues:         xvalues,
		YValues:         yvalues,
		XFormatter:      xformatter,
		YFormatter:      yFormatter,
		TooltipProvider: tooltipProvider,
		Legend:          legend,
	}
}

func NewLineChartView(
	lines []LineView,
	xaxis AxisView,
	yaxis AxisView,
) LineChartView {
	return LineChartView{
		Lines: lines,
		X:     xaxis,
		Y:     yaxis,
	}
}

func TimeFormatter(v float64) string {
	return time.UnixMilli(int64(v)).Format(time.TimeOnly)
}

func BitrateFormatter(v float64) string {
	return models.BitRate(v).String()
}
