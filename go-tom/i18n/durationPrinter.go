package i18n

import (
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
	// " h" is a thin-space character, https://www.compart.com/en/unicode/U+2009
	supportedMapping = map[string]langDef{
		"de": {separator: ":", minSuffix: " h", short: []string{"h", "m", "s"}, long: []string{"Std.", "Min.", "Sek."}},
		"en": {separator: ":", minSuffix: " h", short: []string{"h", "m", "s"}, long: []string{"hrs.", "min", "sec"}},
	}
}

type DurationPrinter interface {
	Minimal(duration time.Duration, showSeconds bool) string
	Short(duration time.Duration, showSeconds bool) string
	Long(duration time.Duration, showSeconds bool) string
}

func NewDurationPrinter(lang language.Tag) DurationPrinter {
	base, _ := lang.Base()
	return &durationPrinter{spec: supportedMapping[base.String()]}
}

func NewDecimalDurationPrinter(lang language.Tag) DurationPrinter {
	base, _ := lang.Base()
	return &decimalDurationPrinter{printer: FindLocale(lang, false), spec: supportedMapping[base.String()]}
}
