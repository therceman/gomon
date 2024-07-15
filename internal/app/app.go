// internal/app/app.go

package app

import (
	"github.com/therceman/gomon/internal/stats"
	"log"
	"time"
)

type Config struct {
	ContainerName      string
	GrafanaInfluxURL   string
	GrafanaAPIKey      string
	GrafanaUsername    string
	ReadTickerTimeSec  int
	FlushTickerTimeSec int
	SendSelfStats      bool
	DockerMetricKeys   []string
}

func Run(config *Config) {
	log.Println("Running Go Monitor for Container")

	readTickerTime := time.Duration(config.ReadTickerTimeSec) * time.Second
	flushTickerTime := time.Duration(config.FlushTickerTimeSec) * time.Second

	ticker := time.NewTicker(readTickerTime)
	flushTicker := time.NewTicker(flushTickerTime)

	defer ticker.Stop()
	defer flushTicker.Stop()

	statsMap := make(map[string]*stats.Stats)

	for {
		select {
		case <-ticker.C:
			systemFetchError := stats.FetchSystemStats(statsMap, config.ContainerName)
			if systemFetchError != nil {
				log.Printf("Error fetching system stats: %v", systemFetchError)
			}
			dockerFetchError := stats.FetchDockerStats(statsMap)
			if dockerFetchError != nil {
				log.Printf("Error fetching docker stats: %v", dockerFetchError)
			}
			workerFetchError := stats.FetchWorkerStats(statsMap, "gomon")
			if workerFetchError != nil {
				log.Printf("Error fetching worker stats: %v", workerFetchError)
			}
		case <-flushTicker.C:
			stats.FlushStats(statsMap)
			statsMap = make(map[string]*stats.Stats)
		}
	}
}
