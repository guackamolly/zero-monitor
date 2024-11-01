package repositories

import (
	"errors"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/showwin/speedtest-go/speedtest"
)

const ooklaServerMediator = "Ookla"

type NetSpeedtestRepository struct {
	client        *speedtest.Speedtest
	closestServer *speedtest.Server
}

func NewNetSpeedtestRepository(
	client *speedtest.Speedtest,
) *NetSpeedtestRepository {
	return &NetSpeedtestRepository{
		client: client,
	}
}

func (r NetSpeedtestRepository) Start() (chan (models.Speedtest), error) {
	var err error
	if r.closestServer == nil {
		err = r.cacheClosestServer()
	}

	if err != nil {
		return nil, err
	}

	ch := make(chan (models.Speedtest))
	go func() {
		// 1. Init phase
		srv := r.closestServer

		st := models.NewSpeedtest(models.UUID(), srv.Sponsor, srv.Name, ooklaServerMediator, srv.Distance*1000)
		ch <- st

		r.client.SetCallbackDownload(func(downRate speedtest.ByteRate) {
			st = st.WithUpdatedDownloadSpeed(float64(downRate) / 0.125)
			println(downRate.String())
			ch <- st
		})

		r.client.SetCallbackUpload(func(upRate speedtest.ByteRate) {
			st = st.WithUpdatedUploadSpeed(float64(upRate) / 0.125)
			ch <- st
		})

		// 2. Latency phase
		st = st.NextPhase()
		ch <- st

		err = srv.PingTest(nil)
		st = st.WithUpdatedLatency(int64(srv.Latency)).NextPhase()
		ch <- st

		// 3. Download phase
		err = srv.DownloadTest()
		st = st.NextPhase()
		ch <- st

		// 4. Upload phase
		err = srv.UploadTest()
		st = st.NextPhase()
		ch <- st

		srv.Context.Reset()
		close(ch)
	}()

	return ch, nil
}

func (r *NetSpeedtestRepository) cacheClosestServer() error {
	srvs, err := r.client.FetchServers()
	if err != nil {
		return err
	}

	srvs, err = srvs.FindServer([]int{})
	if err != nil {
		return err
	}

	if len(srvs) == 0 {
		return errors.New("unexpected error: find closest server did not error, but returned empty list")
	}

	r.closestServer = srvs[0]
	return nil
}
