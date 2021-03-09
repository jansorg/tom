package i18n

import (
	"fmt"
	"time"
)

type isoDelegate struct {
	locateDelegate
}

func (i *isoDelegate) FmtDateShort(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func (i *isoDelegate) FmtDateMedium(t time.Time) string {
	return fmt.Sprintf("%02d %s %04d", t.Day(), i.delegate.MonthAbbreviated(t.Month()), t.Year())
}

func (i *isoDelegate) FmtDateLong(t time.Time) string {
	return fmt.Sprintf("%02d %s %04d", t.Day(), i.delegate.MonthWide(t.Month()), t.Year())
}

func (i *isoDelegate) FmtDateFull(t time.Time) string {
	return fmt.Sprintf("%s, %02d %s %04d", i.delegate.WeekdayWide(t.Weekday()), t.Day(), i.delegate.MonthWide(t.Month()), t.Year())
}
