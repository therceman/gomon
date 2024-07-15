// internal/utils/converter.go

package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// -------- Convert Uint64 to Uint*

// ConvertUint64toUint8 ensures the value fits within the range of uint8.

// ConvertUint64toUint16 ensures the value fits within the range of uint16.

// ConvertUint64toUint32 ensures the value fits within the range of uint32.

// --------- Convert Uint32 to Uint*

// ConvertUint32toUint8 ensures the value fits within the range of uint8.

// ConvertUint32toUint16 ensures the value fits within the range of uint16.

// ConvertUint32toUint64 simply converts the value to uint64.

// --------- Convert String to Uint*

// ConvertStringToUint32 converts a string to uint32, ensuring it fits within the range of uint32.
func ConvertStringToUint32(s string) (uint32, error) {
	value, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(value), nil
}

// ConvertStringToUint64 converts a string to uint64.

// --------- Convert String to Float*

// ConvertStringToFloat32 converts a string to float32.
func ConvertStringToFloat32(s string) (float32, error) {
	value, err := strconv.ParseFloat(s, 32) // Base 10 means the string is interpreted as a decimal number
	if err != nil {
		return 0, err
	}
	return float32(value), nil
}

// ---------- Convert Uint* to String

// ConvertUint32ToString converts uint32 value to a string.
func ConvertUint32ToString(value uint32) string {
	return strconv.FormatUint(uint64(value), 10)
}

// ------------------------------- Rounding -------------------------------

func RoundToTwoDecimal(value float32) float32 {
	return float32(int(value*100)) / 100
}

func ConvertToPerc(value string) (float32, error) {
	trimmed := strings.TrimSuffix(value, "%")
	parsedValue, err := strconv.ParseFloat(trimmed, 32)
	if err != nil {
		return 0, err
	}
	return float32(parsedValue), nil
}

func ConvertMemoryToMB(memUsageStr string) (float32, error) {
	memUsageStr = strings.ToUpper(memUsageStr)

	if strings.HasSuffix(memUsageStr, "GIB") {
		memUsageStr = strings.TrimSuffix(memUsageStr, "GIB")
		memUsage, err := strconv.ParseFloat(memUsageStr, 64)
		if err != nil {
			return 0, err
		}
		return float32(memUsage * 1024), nil
	}

	if strings.HasSuffix(memUsageStr, "MIB") {
		memUsageStr = strings.TrimSuffix(memUsageStr, "MIB")
		memUsage, err := strconv.ParseFloat(memUsageStr, 64)
		if err != nil {
			return 0, err
		}
		return float32(memUsage), nil
	}

	if strings.HasSuffix(memUsageStr, "KIB") {
		memUsageStr = strings.TrimSuffix(memUsageStr, "KIB")
		memUsage, err := strconv.ParseFloat(memUsageStr, 64)
		if err != nil {
			return 0, err
		}
		return float32(memUsage / 1024), nil
	}

	return 0, fmt.Errorf("unknown memory unit")
}

func ConvertSizeToMB(value string) (float32, error) {
	value = strings.ToUpper(value)
	if strings.HasSuffix(value, "GB") {
		num, err := strconv.ParseFloat(strings.TrimSuffix(value, "GB"), 64)
		if err != nil {
			return 0, err
		}
		return float32(num * 1024), nil
	}
	if strings.HasSuffix(value, "MB") {
		num, err := strconv.ParseFloat(strings.TrimSuffix(value, "MB"), 64)
		if err != nil {
			return 0, err
		}
		return float32(num), nil
	}
	if strings.HasSuffix(value, "KB") {
		num, err := strconv.ParseFloat(strings.TrimSuffix(value, "KB"), 64)
		if err != nil {
			return 0, err
		}
		return float32(num / 1024), nil
	}
	if strings.HasSuffix(value, "B") {
		num, err := strconv.ParseFloat(strings.TrimSuffix(value, "B"), 64)
		if err != nil {
			return 0, err
		}
		return float32(num / (1024 * 1024)), nil
	}
	return 0, fmt.Errorf("unknown size unit")
}
