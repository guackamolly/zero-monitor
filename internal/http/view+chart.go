package http

import (
	"fmt"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

// todo...
// Y axis define lista dos markers
// depois o Y dos pontos das linhas tem com base estes markers
// exemplo: Marker: 80 Mbps, {0, 170} e Marker: 100 Mbps, {0, 220}
// download speed 85 Mbps, o ponto sera:
// X: ?
// Y: (85 * 170)/80

type LineChartView struct {
	Xaxis models.Axis
	Yaxis models.Axis
	Lines []models.Line
}

func NewSpeedtestHistoryChart(
	speedtests []models.Speedtest,
) LineChartView {
	min := 0.0
	max := 250.0

	var maxdspeed models.BitRate
	var maxuspeed models.BitRate
	times := make([]time.Time, len(speedtests))

	for i, st := range speedtests {
		if st.DownloadSpeed > maxdspeed {
			maxdspeed = st.DownloadSpeed
		}

		if st.UploadSpeed > maxuspeed {
			maxuspeed = st.UploadSpeed
		}

		times[i] = st.TakenAt
	}

	maxy := maxdspeed
	if maxuspeed > maxy {
		maxy = maxuspeed
	}

	xaxis := models.NewTimeAxis(models.NewPoint(min, max), models.NewPoint(max, max), times)
	yaxis := models.NewAxis(models.NewPoint(min, min), models.NewPoint(min, max), 6, func(i int, x, y float64) models.Marker {
		return models.NewValueMarker(x, y, float64(i), fmt.Sprintf("%0.0f", float64(i*20.0)))
	})

	dmarkers := []models.Marker{}
	umarkers := []models.Marker{}
	for i, st := range speedtests {
		dmarkers = append(
			dmarkers,
			models.NewValueMarker(
				xaxis.Fit(float64(st.TakenAt.UnixMilli())),
				yaxis.Fit(speedtests[i].DownloadSpeed.Value()),
				speedtests[i].DownloadSpeed.Value(),
				speedtests[i].DownloadSpeed.String(),
			),
		)

		umarkers = append(
			umarkers,
			models.NewValueMarker(
				xaxis.Fit(float64(st.TakenAt.UnixMilli())),
				yaxis.Fit(speedtests[i].UploadSpeed.Value()),
				speedtests[i].UploadSpeed.Value(),
				speedtests[i].UploadSpeed.String(),
			),
		)
	}

	for _, m := range dmarkers {
		fmt.Printf("m: %v\n", m)
	}

	return LineChartView{
		Xaxis: xaxis,
		Yaxis: yaxis,
		Lines: []models.Line{
			models.NewLine(dmarkers),
			models.NewLine(umarkers),
		},
	}
}
