package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func ToSeconds(interval string) (int64, error) {
	re := regexp.MustCompile(`^(\d+)([smhdHD])$`)
	matches := re.FindStringSubmatch(interval)
	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid interval format: %s", interval)
	}

	value, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid interval value: %s", matches[1])
	}

	unit := matches[2]
	switch unit {
	case "s":
		return value, nil
	case "m":
		return value * 60, nil
	case "h", "H":
		return value * 60 * 60, nil
	case "d", "D":
		return value * 60 * 60 * 24, nil
	default:
		return 0, fmt.Errorf("invalid interval unit: %s", unit)
	}
}
func ToIntervalString(seconds int64) string {
	if seconds <= 0 {
		return "0s"
	}

	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	} else if seconds < 60*60 {
		return fmt.Sprintf("%dm", seconds/60)
	} else if seconds < 60*60*24 {
		return fmt.Sprintf("%dh", seconds/(60*60))
	} else {
		return fmt.Sprintf("%dd", seconds/(60*60*24))
	}
}
func ToMinutes() {}
