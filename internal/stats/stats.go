// internal/stats/stats.go

package stats

import (
	"github.com/therceman/gomon/internal/stats/docker"
	"github.com/therceman/gomon/internal/stats/system"
	"github.com/therceman/gomon/internal/stats/worker"
	"github.com/therceman/gomon/internal/utils"
	"log"
)

type Stats struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	CPUMaxPerc   float32 `json:"cpu_max_perc"`
	CPUMinPerc   float32 `json:"cpu_min_perc"`
	CPUAvgPerc   float32 `json:"cpu_avg_perc"`
	CPUPercSum   float32 `json:"-"` // Used for calculating average
	CPUCount     int     `json:"-"` // Used for calculating average
	MemMaxMB     float32 `json:"mem_max_mb"`
	MemMinMB     float32 `json:"mem_min_mb"`
	MemAvgMB     float32 `json:"mem_avg_mb"`
	MemMBPercSum float32 `json:"-"` // Used for calculating average
	MemCount     int     `json:"-"` // Used for calculating average
	MemMaxPerc   float32 `json:"mem_max_perc"`
	MemMinPerc   float32 `json:"mem_min_perc"`
	MemAvgPerc   float32 `json:"mem_avg_perc"`
	MemPercSum   float32 `json:"-"` // Used for calculating average
	MemPercCount int     `json:"-"` // Used for calculating average
	DiskMB       float32 `json:"disk_mb"`
}

func FlushStats(statsMap map[string]*Stats) {
	log.Println("Flushing stats map")

	for _, value := range statsMap {
		// Log only the desired fields
		log.Printf(
			"ID: %s, Name: %s, CPUMaxPerc: %.2f, CPUAveragePerc: %.2f, "+
				"MemMaxMB: %.2f, MemAverageMB: %.2f, MemMaxPerc: %.2f, MemAveragePerc: %.2f, DiskMB: %.2f\n",
			value.ID, value.Name, value.CPUMaxPerc, value.CPUAvgPerc, value.MemMaxMB, value.MemAvgMB,
			value.MemMaxPerc, value.MemAvgPerc, value.DiskMB)
	}
}

// FetchDockerStats fetches and updates Docker stats
func FetchDockerStats(statsMap map[string]*Stats) error {
	dockerStats, err := docker.GetStats()
	if err != nil {
		return err
	}

	for _, stat := range dockerStats {
		if existing, found := statsMap[stat.ID]; found {
			// Update CPUPerc percentages
			if stat.CPU < existing.CPUMinPerc {
				existing.CPUMinPerc = stat.CPU
			}
			if stat.CPU > existing.CPUMaxPerc {
				existing.CPUMaxPerc = stat.CPU
			}

			// Update CPU average
			existing.CPUPercSum += stat.CPU
			existing.CPUCount++
			existing.CPUAvgPerc = utils.RoundToTwoDecimal(existing.CPUPercSum / float32(existing.CPUCount))

			// Update Memory usage in MB
			if stat.MemMB < existing.MemMinMB {
				existing.MemMinMB = stat.MemMB
			}
			if stat.MemMB > existing.MemMaxMB {
				existing.MemMaxMB = stat.MemMB
			}

			// Update Memory average
			existing.MemMBPercSum += stat.MemMB
			existing.MemCount++
			existing.MemAvgMB = utils.RoundToTwoDecimal(existing.MemMBPercSum / float32(existing.MemCount))

			// Update Memory percentages
			if stat.MemPerc < existing.MemMinPerc {
				existing.MemMinPerc = stat.MemPerc
			}
			if stat.MemPerc > existing.MemMaxPerc {
				existing.MemMaxPerc = stat.MemPerc
			}

			// Update Memory percentage average
			existing.MemPercSum += stat.MemPerc
			existing.MemPercCount++
			existing.MemAvgPerc = utils.RoundToTwoDecimal(existing.MemPercSum / float32(existing.MemPercCount))

			// Update Disk usage
			existing.DiskMB = stat.SizeMB
		} else {
			statsMap[stat.ID] = &Stats{
				ID:           stat.ID,
				Name:         stat.Name,
				CPUMinPerc:   stat.CPU,
				CPUMaxPerc:   stat.CPU,
				CPUPercSum:   stat.CPU,
				CPUCount:     1,
				CPUAvgPerc:   stat.CPU,
				MemMinMB:     stat.MemMB,
				MemMaxMB:     stat.MemMB,
				MemMBPercSum: stat.MemMB,
				MemCount:     1,
				MemAvgMB:     stat.MemMB,
				MemMinPerc:   stat.MemPerc,
				MemMaxPerc:   stat.MemPerc,
				MemPercSum:   stat.MemPerc,
				MemPercCount: 1,
				MemAvgPerc:   stat.MemPerc,
				DiskMB:       stat.SizeMB,
			}
		}
	}

	return nil
}

// FetchSystemStats fetches and updates system stats
func FetchSystemStats(statsMap map[string]*Stats, containerName string) error {
	sysStats, err := system.GetStats()
	if err != nil {
		return err
	}

	ID := "system"
	NAME := containerName

	if existing, found := statsMap[ID]; found {
		// Update CPUPerc percentages
		if sysStats.CPUPerc < existing.CPUMinPerc {
			existing.CPUMinPerc = sysStats.CPUPerc
		}
		if sysStats.CPUPerc > existing.CPUMaxPerc {
			existing.CPUMaxPerc = sysStats.CPUPerc
		}

		// Update CPU average
		existing.CPUPercSum += sysStats.CPUPerc
		existing.CPUCount++
		existing.CPUAvgPerc = utils.RoundToTwoDecimal(existing.CPUPercSum / float32(existing.CPUCount))

		// Update Memory usage in MB
		if float32(sysStats.MemMB) < existing.MemMinMB {
			existing.MemMinMB = float32(sysStats.MemMB)
		}
		if float32(sysStats.MemMB) > existing.MemMaxMB {
			existing.MemMaxMB = float32(sysStats.MemMB)
		}

		// Update Memory average
		existing.MemMBPercSum += float32(sysStats.MemMB)
		existing.MemCount++
		existing.MemAvgMB = utils.RoundToTwoDecimal(existing.MemMBPercSum / float32(existing.MemCount))

		// Update Memory percentages
		if sysStats.MemPerc < existing.MemMinPerc {
			existing.MemMinPerc = sysStats.MemPerc
		}
		if sysStats.MemPerc > existing.MemMaxPerc {
			existing.MemMaxPerc = sysStats.MemPerc
		}

		// Update Memory percentage average
		existing.MemPercSum += sysStats.MemPerc
		existing.MemPercCount++
		existing.MemAvgPerc = utils.RoundToTwoDecimal(existing.MemPercSum / float32(existing.MemPercCount))

		// Update Disk usage
		existing.DiskMB = float32(sysStats.DiskMB)
	} else {
		statsMap[ID] = &Stats{
			ID:           ID,
			Name:         NAME,
			CPUMinPerc:   sysStats.CPUPerc,
			CPUMaxPerc:   sysStats.CPUPerc,
			CPUPercSum:   sysStats.CPUPerc,
			CPUCount:     1,
			CPUAvgPerc:   sysStats.CPUPerc,
			MemMinMB:     float32(sysStats.MemMB),
			MemMaxMB:     float32(sysStats.MemMB),
			MemMBPercSum: float32(sysStats.MemMB),
			MemCount:     1,
			MemAvgMB:     float32(sysStats.MemMB),
			MemMinPerc:   sysStats.MemPerc,
			MemMaxPerc:   sysStats.MemPerc,
			MemPercSum:   sysStats.MemPerc,
			MemPercCount: 1,
			MemAvgPerc:   sysStats.MemPerc,
			DiskMB:       float32(sysStats.DiskMB),
		}
	}

	return nil
}

// FetchWorkerStats fetches and updates worker stats
func FetchWorkerStats(statsMap map[string]*Stats, processName string) error {
	workerStats, err := worker.GetStats(processName)
	if err != nil {
		return err
	}

	ID := utils.ConvertUint32ToString(workerStats.PID)
	NAME := processName

	if existing, found := statsMap[ID]; found {
		// Update CPU percentages
		if workerStats.CPUPerc < existing.CPUMinPerc {
			existing.CPUMinPerc = workerStats.CPUPerc
		}
		if workerStats.CPUPerc > existing.CPUMaxPerc {
			existing.CPUMaxPerc = workerStats.CPUPerc
		}

		// Update CPU average
		existing.CPUPercSum += workerStats.CPUPerc
		existing.CPUCount++
		existing.CPUAvgPerc = utils.RoundToTwoDecimal(existing.CPUPercSum / float32(existing.CPUCount))

		// Update Memory usage in MB
		memMB := utils.RoundToTwoDecimal(float32(workerStats.MemKB) / 1024)
		if memMB < existing.MemMinMB {
			existing.MemMinMB = memMB
		}
		if memMB > existing.MemMaxMB {
			existing.MemMaxMB = memMB
		}

		// Update Memory average
		existing.MemMBPercSum += memMB
		existing.MemCount++
		existing.MemAvgMB = utils.RoundToTwoDecimal(existing.MemMBPercSum / float32(existing.MemCount))

		// Update Memory percentages
		if workerStats.MemPerc < existing.MemMinPerc {
			existing.MemMinPerc = workerStats.MemPerc
		}
		if workerStats.MemPerc > existing.MemMaxPerc {
			existing.MemMaxPerc = workerStats.MemPerc
		}

		// Update Memory percentage average
		existing.MemPercSum += workerStats.MemPerc
		existing.MemPercCount++
		existing.MemAvgPerc = utils.RoundToTwoDecimal(existing.MemPercSum / float32(existing.MemPercCount))
	} else {
		statsMap[ID] = &Stats{
			ID:           ID,
			Name:         NAME,
			CPUMinPerc:   workerStats.CPUPerc,
			CPUMaxPerc:   workerStats.CPUPerc,
			CPUPercSum:   workerStats.CPUPerc,
			CPUCount:     1,
			CPUAvgPerc:   workerStats.CPUPerc,
			MemMinMB:     utils.RoundToTwoDecimal(float32(workerStats.MemKB) / 1024),
			MemMaxMB:     utils.RoundToTwoDecimal(float32(workerStats.MemKB) / 1024),
			MemMBPercSum: utils.RoundToTwoDecimal(float32(workerStats.MemKB) / 1024),
			MemCount:     1,
			MemAvgMB:     utils.RoundToTwoDecimal(float32(workerStats.MemKB) / 1024),
			MemMinPerc:   workerStats.MemPerc,
			MemMaxPerc:   workerStats.MemPerc,
			MemPercSum:   workerStats.MemPerc,
			MemPercCount: 1,
			MemAvgPerc:   workerStats.MemPerc,
			DiskMB:       0,
		}
	}

	return nil
}
