package system

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/therceman/gomon/internal/utils"
)

// Stats holds combined metrics for system resources
type Stats struct {
	Mem      uint64  `json:"mem"`       // Used memory in MB
	MemPerc  float64 `json:"mem_perc"`  // Used memory percentage
	CPU      float64 `json:"cpu"`       // CPU usage percentage
	Disk     uint64  `json:"disk"`      // Used disk space in MB
	DiskPerc float64 `json:"disk_perc"` // Used disk space percentage
}

// GetStats retrieves system statistics including memory, CPU, and disk usage
func GetStats() (Stats, error) {
	memStats, err := getMemStats()
	if err != nil {
		return Stats{}, err
	}

	cpuStats, err := getCPUStats()
	if err != nil {
		return Stats{}, err
	}

	diskStats, err := getDiskStats("/")
	if err != nil {
		return Stats{}, err
	}

	return Stats{
		Mem:      memStats.Used,
		MemPerc:  memStats.UsedPercent,
		CPU:      cpuStats.CPUUsagePercent,
		Disk:     diskStats.Used,
		DiskPerc: diskStats.UsedPerc,
	}, nil
}

type memStats struct {
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type cpuStats struct {
	CPUUsagePercent float64 `json:"cpu_usage_percent"`
}

type diskStats struct {
	Used     uint64  `json:"used"`
	UsedPerc float64 `json:"used_perc"`
}

func getMemStats() (memStats, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return memStats{}, err
	}
	defer file.Close()

	var total, free uint64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			total, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return memStats{}, err
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			free, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return memStats{}, err
			}
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return memStats{}, err
	}

	used := total - free
	usedPercent := (float64(used) / float64(total)) * 100

	return memStats{
		Used:        used / 1024, // Convert to MB
		UsedPercent: utils.RoundToTwoDecimal(usedPercent),
	}, nil
}

// getCPUStats retrieves CPU usage percentage
func getCPUStats() (cpuStats, error) {
	idle0, total0, _, _ := getCPUSample()
	// Sleep for a short duration to get a new sample
	time.Sleep(200 * time.Millisecond)
	idle1, total1, _, _ := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)

	// Avoid division by zero
	if totalTicks == 0 {
		return cpuStats{CPUUsagePercent: 0}, nil
	}

	cpuUsage := (1.0 - (idleTicks / totalTicks)) * 100

	return cpuStats{
		CPUUsagePercent: utils.RoundToTwoDecimal(cpuUsage),
	}, nil
}

// getCPUSample reads CPU usage statistics from /proc/stat
// and returns per-core idle and total times
func getCPUSample() (idle, total uint64, idleTimes, totalTimes []uint64) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, nil, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) > 0 && strings.HasPrefix(fields[0], "cpu") {
			idleTime, totalTime := uint64(0), uint64(0)
			for i, v := range fields[1:] {
				val, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return 0, 0, nil, nil
				}
				totalTime += val
				if i == 3 { // idle time is the 4th field
					idleTime = val
				}
			}
			if fields[0] == "cpu" {
				total += totalTime
				idle += idleTime
			} else {
				idleTimes = append(idleTimes, idleTime)
				totalTimes = append(totalTimes, totalTime)
			}
		}
	}
	return idle, total, idleTimes, totalTimes
}

func getDiskStats(path string) (diskStats, error) {
	var stat syscall.Statfs_t

	if err := syscall.Statfs(path, &stat); err != nil {
		return diskStats{}, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := total - free
	usedPerc := (float64(used) / float64(total)) * 100

	return diskStats{
		Used:     used / 1024 / 1024, // Convert to MB
		UsedPerc: utils.RoundToTwoDecimal(usedPerc),
	}, nil
}
