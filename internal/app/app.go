package app

import (
	"log"

	"github.com/therceman/gomon/internal/stats/docker"
	"github.com/therceman/gomon/internal/stats/system"
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

	// test output
	// TODO: real-time sync

	systemStats, err := system.GetStats()
	if err != nil {
		log.Fatalf("Error getting system stats: %v", err)
	}
	log.Printf("System Stats: %+v\n", systemStats)

	dockerStats, err := docker.GetStats()
	if err != nil {
		log.Fatalf("Error getting system stats: %v", err)
	}
	log.Printf("Docker Stats: %+v\n", dockerStats)
}
