package i18n

import (
	"fmt"
	"time"

	"github.com/go-playground/locales"
)

type decimalDurationPrinter struct {
	printer locales.Translator
	spec    langDef
}

func (p *decimalDurationPrinter) Minimal(d time.Duration, showSeconds bool) string {
	s := d.Seconds()
	hours := float64(s) / 3600.0
	return fmt.Sprintf("%s%s", p.printer.FmtNumber(hours, 2), p.spec.minSuffix)
}

func (p *decimalDurationPrinter) Short(d time.Duration, showSeconds bool) string {
	return p.Minimal(d, showSeconds)
}

func (p *decimalDurationPrinter) Long(d time.Duration, showSeconds bool) string {
	return p.Minimal(d, showSeconds)
}
