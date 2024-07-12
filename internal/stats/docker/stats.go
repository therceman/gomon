package docker

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/therceman/gomon/internal/utils"
)

type DockerStats struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	CPU     float64 `json:"cpu"`
	Mem     float64 `json:"mem"`
	MemPerc float64 `json:"mem_perc"`
	NetI    float64 `json:"net_i"`
	NetO    float64 `json:"net_o"`
	BlockI  float64 `json:"block_i"`
	BlockO  float64 `json:"block_o"`
	PIDs    int     `json:"pids"`
	Size    float64 `json:"size"`
}

func parseDockerStatsOutput(output string) ([]DockerStats, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var stats []DockerStats
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "CONTAINER ID") {
			continue // Skip the header line
		}

		fields := strings.Fields(line)
		if len(fields) < 14 {
			log.Printf("Skipping incomplete line: %s", line)
			continue // Skip incomplete lines
		}

		netI, err := utils.ConvertSizeToMB(fields[7])
		if err != nil {
			log.Printf("Error converting Net I: %v", err)
			continue
		}

		netO, err := utils.ConvertSizeToMB(fields[9])
		if err != nil {
			log.Printf("Error converting Net O: %v", err)
			continue
		}

		blockI, err := utils.ConvertSizeToMB(fields[10])
		if err != nil {
			log.Printf("Error converting Block I: %v", err)
			continue
		}

		blockO, err := utils.ConvertSizeToMB(fields[12])
		if err != nil {
			log.Printf("Error converting Block O: %v", err)
			continue
		}

		cpuUsage, err := utils.ConvertToPerc(fields[2])
		if err != nil {
			log.Printf("Error parsing CPU usage: %v", err)
			continue
		}

		memUsage, err := utils.ConvertMemoryToMB(fields[3])
		if err != nil {
			log.Printf("Error parsing memory usage: %v", err)
			continue
		}

		memPerc, err := utils.ConvertToPerc(fields[6])
		if err != nil {
			log.Printf("Error parsing memory percent: %v", err)
			continue
		}

		pids, err := strconv.Atoi(fields[13])
		if err != nil {
			log.Printf("Error parsing PIDs: %v", err)
			continue
		}

		containerSize, err := GetContainerSize(fields[0], true)
		if err != nil {
			log.Printf("Error getting docker container size: %v", err)
			continue
		}

		stat := DockerStats{
			ID:      fields[0],
			Name:    fields[1],
			CPU:     utils.RoundToTwoDecimal(cpuUsage),
			Mem:     utils.RoundToTwoDecimal(memUsage),
			MemPerc: utils.RoundToTwoDecimal(memPerc),
			NetI:    utils.RoundToTwoDecimal(netI),
			NetO:    utils.RoundToTwoDecimal(netO),
			BlockI:  utils.RoundToTwoDecimal(blockI),
			BlockO:  utils.RoundToTwoDecimal(blockO),
			PIDs:    pids,
			Size:    utils.RoundToTwoDecimal(containerSize),
		}
		stats = append(stats, stat)
	}
	return stats, nil
}

func GetStats() ([]DockerStats, error) {
	cmd := exec.Command("docker", "stats", "--no-stream")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	output := out.String()

	stats, err := parseDockerStatsOutput(output)
	if err != nil {
		return nil, err
	}

	// Ensure we do not hold on to memory longer than needed
	out.Reset()
	return stats, nil
}
