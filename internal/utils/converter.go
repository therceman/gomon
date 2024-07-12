package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func RoundToTwoDecimal(value float64) float64 {
	return float64(int(value*100)) / 100
}

func ConvertToPerc(value string) (float64, error) {
	trimmed := strings.TrimSuffix(value, "%")

	return strconv.ParseFloat(trimmed, 64)
}

func ConvertMemoryToMB(memUsageStr string) (float64, error) {
	memUsageStr = strings.ToUpper(memUsageStr)

	if strings.HasSuffix(memUsageStr, "GIB") {
		memUsageStr = strings.TrimSuffix(memUsageStr, "GIB")
		memUsage, err := strconv.ParseFloat(memUsageStr, 64)
		if err != nil {
			return 0, err
		}
		return memUsage * 1024, nil
	}

	if strings.HasSuffix(memUsageStr, "MIB") {
		memUsageStr = strings.TrimSuffix(memUsageStr, "MIB")
		return strconv.ParseFloat(memUsageStr, 64)
	}

	if strings.HasSuffix(memUsageStr, "KIB") {
		memUsageStr = strings.TrimSuffix(memUsageStr, "KIB")
		memUsage, err := strconv.ParseFloat(memUsageStr, 64)
		if err != nil {
			return 0, err
		}
		return memUsage / 1024, nil
	}

	return 0, fmt.Errorf("unknown memory unit")
}

func ConvertSizeToMB(value string) (float64, error) {
	value = strings.ToUpper(value)
	if strings.HasSuffix(value, "GB") {
		num, err := strconv.ParseFloat(strings.TrimSuffix(value, "GB"), 64)
		if err != nil {
			return 0, err
		}
		return num * 1024, nil
	}
	if strings.HasSuffix(value, "MB") {
		num, err := strconv.ParseFloat(strings.TrimSuffix(value, "MB"), 64)
		if err != nil {
			return 0, err
		}
		return num, nil
	}
	if strings.HasSuffix(value, "KB") {
		num, err := strconv.ParseFloat(strings.TrimSuffix(value, "KB"), 64)
		if err != nil {
			return 0, err
		}
		return num / 1024, nil
	}
	if strings.HasSuffix(value, "B") {
		num, err := strconv.ParseFloat(strings.TrimSuffix(value, "B"), 64)
		if err != nil {
			return 0, err
		}
		return num / (1024 * 1024), nil
	}
	return 0, fmt.Errorf("unknown size unit")
}
