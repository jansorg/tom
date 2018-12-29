package i18n

import (
	"fmt"
	"time"

	"golang.org/x/text/language"
)

type langDef struct {
	separator string
	minSuffix string
	short     []string
	long      []string
}

var supportedMapping map[string]langDef

func init() {
	supportedMapping = map[string]langDef{
		"de": {separator: ":", minSuffix: " h", short: []string{"h", "m", "s"}, long: []string{"Std.", "Min.", "Sek."}},
		"en": {separator: ":", short: []string{"h", "m", "s"}, long: []string{"hrs.", "min", "sec"}},
	}
}

type DurationPrinter interface {
	Minimal(duration time.Duration) string
	Short(duration time.Duration) string
	Long(duration time.Duration) string
}

func NewDurationPrinter(lang language.Tag) DurationPrinter {
	base, _ := lang.Base()
	return &durationPrinter{spec: supportedMapping[base.String()]}
}

type durationPrinter struct {
	spec langDef
}

func (p *durationPrinter) Minimal(d time.Duration) string {
	s := d.Seconds()
	hours := int(s / 3600.0)
	minutes := int(int(s) % 3600 / 60.0)
	seconds := int(s) % 3600 % 60

	return fmt.Sprintf("%d%s%02d%s%02d%s", hours, p.spec.separator, minutes, p.spec.separator, seconds, p.spec.minSuffix)
}

func (p *durationPrinter) Short(d time.Duration) string {
	return fmt.Sprintf("%d%s%d%s%d%s", d.Hours(), p.spec.short[0], d.Minutes(), p.spec.short[1], d.Seconds(), p.spec.short[2])
}

func (p *durationPrinter) Long(d time.Duration) string {
	return fmt.Sprintf("%d%s%d%s%d%s", d.Hours(), p.spec.long[0], d.Minutes(), p.spec.long[1], d.Seconds(), p.spec.long[2])
}
