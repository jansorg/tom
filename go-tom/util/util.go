package util

import (
	"fmt"
	"log"
	"time"
)



func ParseRoundingMode(mode string) RoundingMode {
	if mode == "" {
		return RoundNone
	}

	switch mode {
	case "up":
		return RoundUp
	case "nearest":
		return RoundNearest
	default:
		log.Fatal(fmt.Errorf("unknown rounding mode %s. Possible values: up, nearest", mode))
		return RoundNone
	}
}

func RoundDuration(value time.Duration, config RoundingConfig) time.Duration {
	if config.Mode == RoundNone {
		return value
	}

	result := value.Round(config.Size)
	if config.Mode == RoundNearest {
		return result
	}

	if result < value {
		// round up if Go rounded down
		result = result + config.Size
	}
	return result
}

func ShortDateString(date time.Time) string {
	return date.Format("2006-01-02")
}
