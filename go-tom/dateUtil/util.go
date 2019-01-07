package dateUtil

import (
	"fmt"
	"log"
	"time"
)

type RoundingMode int8

const (
	RoundNone RoundingMode = iota + 1
	RoundNearest
	RoundUp
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

func RoundDuration(value time.Duration, mode RoundingMode, roundTo time.Duration) time.Duration {
	if mode == RoundNone {
		return value
	}

	result := value.Round(roundTo)
	if mode == RoundNearest {
		return result
	}

	if result < value {
		// round up if Go rounded down
		result = result + roundTo
	}
	return result
}

func ShortDateString(date time.Time) string {
	return date.Format("2006-01-02")
}