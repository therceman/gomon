package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/therceman/gomon/internal/app"
	"github.com/therceman/gomon/internal/dotenv"
)

// LoadConfig loads environment variables into a Config struct
func LoadConfig() (*app.Config, error) {
	readTickerTimeSec, err := strconv.Atoi(os.Getenv("READ_TICKER_TIME_SEC"))
	if err != nil {
		return nil, fmt.Errorf("invalid value for READ_TICKER_TIME_SEC")
	}

	flushTickerTimeSec, err := strconv.Atoi(os.Getenv("FLUSH_TICKER_TIME_SEC"))
	if err != nil {
		return nil, fmt.Errorf("invalid value for FLUSH_TICKER_TIME_SEC")
	}

	sendSelfStats, err := strconv.ParseBool(os.Getenv("SEND_SELF_STATS"))
	if err != nil {
		return nil, fmt.Errorf("invalid value for SEND_SELF_STATS")
	}

	var dockerMetricKeys []string
	keys := os.Getenv("DOCKER_METRIC_KEYS")
	if keys == "" {
		dockerMetricKeys = []string{"cpu_max", "mem_max", "net_i", "net_o", "block_i", "block_o"}
	} else {
		dockerMetricKeys = strings.Split(keys, ",")
	}

	config := &app.Config{
		ContainerName:      os.Getenv("CONTAINER_NAME"),
		GrafanaInfluxURL:   os.Getenv("GRAFANA_INFLUX_URL"),
		GrafanaAPIKey:      os.Getenv("GRAFANA_API_KEY"),
		GrafanaUsername:    os.Getenv("GRAFANA_USERNAME"),
		ReadTickerTimeSec:  readTickerTimeSec,
		FlushTickerTimeSec: flushTickerTimeSec,
		SendSelfStats:      sendSelfStats,
		DockerMetricKeys:   dockerMetricKeys,
	}

	return config, nil
}

func main() {
	// Load the environment variables from the .env file
	err := dotenv.LoadEnv(".env")
	if err != nil {
		log.Fatalf("Could not load .env file: %v", err)
	}

	// Load the configuration into the Config struct
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Pass the Config struct to the app
	app.Run(config)
}
