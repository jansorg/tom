package i18n

import (
	"fmt"
	"time"
)

type durationPrinter struct {
	spec langDef
}

func (p *durationPrinter) Minimal(d time.Duration, showSeconds bool) string {
	s := d.Seconds()
	hours := int(s / 3600.0)
	minutes := int(int(s) % 3600 / 60.0)

	if showSeconds {
		seconds := int(s) % 3600 % 60
		return fmt.Sprintf("%d%s%02d%s%02d%s", hours, p.spec.separator, minutes, p.spec.separator, seconds, p.spec.minSuffix)
	}

	return fmt.Sprintf("%d%s%02d%s", hours, p.spec.separator, minutes, p.spec.minSuffix)
}

func (p *durationPrinter) Short(d time.Duration, showSeconds bool) string {
	s := d.Seconds()
	hours := int(s / 3600.0)
	minutes := int(s) % 3600 / 60.0
	seconds := int(s) % 3600 % 60

	if showSeconds {
		return fmt.Sprintf("%d%s %d%s %d%s", hours, p.spec.short[0], minutes, p.spec.short[1], seconds, p.spec.short[2])
	}

	return fmt.Sprintf("%d%s %d%s", hours, p.spec.short[0], minutes, p.spec.short[1])
}

func (p *durationPrinter) Long(d time.Duration, showSeconds bool) string {
	s := d.Seconds()
	hours := int(s / 3600.0)
	minutes := int(s) % 3600 / 60.0
	seconds := int(s) % 3600 % 60

	if showSeconds {
		return fmt.Sprintf("%d%s %d%s %d%s", hours, p.spec.long[0], minutes, p.spec.long[1], seconds, p.spec.long[2])
	}

	return fmt.Sprintf("%d%s %d%s", hours, p.spec.long[0], minutes, p.spec.long[1])
}
