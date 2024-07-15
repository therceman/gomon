// internal/stats/worker/stats.go

package worker

import (
	"bufio"
	"fmt"
	"github.com/therceman/gomon/internal/utils"
	"os/exec"
	"strings"
)

type Stats struct {
	MemKB   uint32  `json:"mem_kb"`   // Used memory in KB
	CPUPerc float32 `json:"cpu_perc"` // CPU usage percentage
	MemPerc float32 `json:"mem_perc"` // Memory usage percentage
	PID     uint32  `json:"pid"`      // Process ID
}

func GetStats(processName string) (Stats, error) {
	psCmd := exec.Command("sh", "-c", fmt.Sprintf(`ps aux | grep "%s" | grep -v grep | awk '{print $2 " " $3 " " $4 " " $6}'`, processName))
	output, err := psCmd.Output()
	if err != nil {
		return Stats{}, fmt.Errorf("error executing ps command: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			return Stats{}, fmt.Errorf("unexpected format in ps output")
		}

		pid, err := utils.ConvertStringToUint32(fields[0])
		if err != nil {
			return Stats{}, err
		}

		cpuPerc, err := utils.ConvertStringToFloat32(fields[1])
		if err != nil {
			return Stats{}, err
		}

		memPerc, err := utils.ConvertStringToFloat32(fields[2])
		if err != nil {
			return Stats{}, err
		}

		memKB, err := utils.ConvertStringToUint32(fields[3])
		if err != nil {
			return Stats{}, err
		}

		return Stats{
			PID:     pid,
			CPUPerc: cpuPerc,
			MemPerc: memPerc,
			MemKB:   memKB,
		}, nil
	}

	if err := scanner.Err(); err != nil {
		return Stats{}, err
	}

	return Stats{}, fmt.Errorf("failed to read process stats")
}
