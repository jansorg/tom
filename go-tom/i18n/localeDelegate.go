package i18n

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/currency"
	"time"
)

type locateDelegate struct {
	delegate locales.Translator
}

func (i *locateDelegate) Locale() string {
	return i.delegate.Locale()
}

func (i *locateDelegate) PluralsCardinal() []locales.PluralRule {
	return i.delegate.PluralsCardinal()
}

func (i *locateDelegate) PluralsOrdinal() []locales.PluralRule {
	return i.delegate.PluralsOrdinal()
}

func (i *locateDelegate) PluralsRange() []locales.PluralRule {
	return i.delegate.PluralsRange()
}

func (i *locateDelegate) CardinalPluralRule(num float64, v uint64) locales.PluralRule {
	return i.delegate.CardinalPluralRule(num, v)
}

func (i *locateDelegate) OrdinalPluralRule(num float64, v uint64) locales.PluralRule {
	return i.delegate.OrdinalPluralRule(num, v)
}

func (i *locateDelegate) RangePluralRule(num1 float64, v1 uint64, num2 float64, v2 uint64) locales.PluralRule {
	return i.delegate.RangePluralRule(num1, v1, num2, v2)
}

func (i *locateDelegate) MonthAbbreviated(month time.Month) string {
	return i.delegate.MonthAbbreviated(month)
}

func (i *locateDelegate) MonthsAbbreviated() []string {
	return i.delegate.MonthsAbbreviated()
}

func (i *locateDelegate) MonthNarrow(month time.Month) string {
	return i.delegate.MonthNarrow(month)
}

func (i *locateDelegate) MonthsNarrow() []string {
	return i.delegate.MonthsNarrow()
}

func (i *locateDelegate) MonthWide(month time.Month) string {
	return i.delegate.MonthWide(month)
}

func (i *locateDelegate) MonthsWide() []string {
	return i.delegate.MonthsWide()
}

func (i *locateDelegate) WeekdayAbbreviated(weekday time.Weekday) string {
	return i.delegate.WeekdayAbbreviated(weekday)
}

func (i *locateDelegate) WeekdaysAbbreviated() []string {
	return i.delegate.WeekdaysAbbreviated()
}

func (i *locateDelegate) WeekdayNarrow(weekday time.Weekday) string {
	return i.delegate.WeekdayNarrow(weekday)
}

func (i *locateDelegate) WeekdaysNarrow() []string {
	return i.delegate.WeekdaysNarrow()
}

func (i *locateDelegate) WeekdayShort(weekday time.Weekday) string {
	return i.delegate.WeekdayShort(weekday)
}

func (i *locateDelegate) WeekdaysShort() []string {
	return i.delegate.WeekdaysShort()
}

func (i *locateDelegate) WeekdayWide(weekday time.Weekday) string {
	return i.delegate.WeekdayWide(weekday)
}

func (i *locateDelegate) WeekdaysWide() []string {
	return i.delegate.WeekdaysWide()
}

func (i *locateDelegate) FmtNumber(num float64, v uint64) string {
	return i.delegate.FmtNumber(num, v)
}

func (i *locateDelegate) FmtPercent(num float64, v uint64) string {
	return i.delegate.FmtPercent(num, v)
}

func (i *locateDelegate) FmtCurrency(num float64, v uint64, currency currency.Type) string {
	return i.delegate.FmtCurrency(num, v, currency)
}

func (i *locateDelegate) FmtAccounting(num float64, v uint64, currency currency.Type) string {
	return i.delegate.FmtAccounting(num, v, currency)
}

func (i *locateDelegate) FmtDateShort(t time.Time) string {
	return i.delegate.FmtDateShort(t)
}

func (i *locateDelegate) FmtDateMedium(t time.Time) string {
	return i.delegate.FmtDateMedium(t)
}

func (i *locateDelegate) FmtDateLong(t time.Time) string {
	return i.delegate.FmtDateLong(t)
}

func (i *locateDelegate) FmtDateFull(t time.Time) string {
	return i.delegate.FmtDateFull(t)
}

func (i *locateDelegate) FmtTimeShort(t time.Time) string {
	return i.delegate.FmtTimeShort(t)
}

func (i *locateDelegate) FmtTimeMedium(t time.Time) string {
	return i.delegate.FmtTimeMedium(t)
}

func (i *locateDelegate) FmtTimeLong(t time.Time) string {
	return i.delegate.FmtTimeLong(t)
}

func (i *locateDelegate) FmtTimeFull(t time.Time) string {
	return i.delegate.FmtTimeFull(t)
}
