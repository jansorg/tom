package dateTime

import (
	"fmt"
	"time"

	"github.com/go-playground/locales"
)

func NewDateRange(start *time.Time, end *time.Time, locale locales.Translator) DateRange {
	if locale == nil {
		panic("locale not set")
	}

	dateRange := DateRange{
		Start:  start,
		End:    end,
		locale: locale,
	}
	dateRange.debug = dateRange.String()
	return dateRange
}

func NewYearRange(date time.Time, locale locales.Translator, location *time.Location) DateRange {
	start := time.Date(date.In(location).Year(), time.January, 1, 0, 0, 0, 0, location)
	end := start.AddDate(1, 0, 0)
	return NewDateRange(&start, &end, locale)
}

func NewMonthRange(date time.Time, locale locales.Translator, location *time.Location) DateRange {
	y, m, _ := date.In(location).Date()
	start := time.Date(y, m, 1, 0, 0, 0, 0, location)
	end := start.AddDate(0, 1, 0)

	return NewDateRange(&start, &end, locale)
}

func NewWeekRange(date time.Time, locale locales.Translator, location *time.Location) DateRange {
	d := date.In(location)

	// e.g. Tuesday = 2-0 -> 2 days to shift back
	daysShift := int(d.Weekday() - time.Sunday)

	start := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	start = start.AddDate(0, 0, -daysShift)

	end := start.AddDate(0, 0, 7)

	return NewDateRange(&start, &end, locale)
}

func NewDayRange(date time.Time, locale locales.Translator, location *time.Location) DateRange {
	y, m, d := date.In(location).Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, location)
	end := start.AddDate(0, 0, 1)

	return NewDateRange(&start, &end, locale)
}

type DateRange struct {
	Start  *time.Time `json:"start"`
	End    *time.Time `json:"end"`
	debug  string
	locale locales.Translator
}

func (r DateRange) In(location *time.Location) DateRange {
	start := r.Start
	if start != nil {
		in := start.In(location)
		start = &in
	}

	end := r.End
	if end != nil {
		in := end.In(location)
		end = &in
	}

	return NewDateRange(start, end, r.locale)
}

func (r DateRange) String() string {
	var start, end string
	if r.Start != nil {
		start = r.Start.String()
	}
	if r.End != nil {
		end = r.End.String()
	}
	return fmt.Sprintf("%s - %s", start, end)
}

func (r DateRange) ShortString() string {
	var start, end string
	if r.Start != nil {
		start = ShortDateString(*r.Start)
	}
	if r.End != nil {
		end = ShortDateString(*r.End)
	}
	return fmt.Sprintf("%s - %s", start, end)
}

func (r DateRange) MinimalString() string {
	// year
	if r.Start.AddDate(1, 0, 0).Equal(*r.End) {
		return fmt.Sprintf("%04d", r.Start.Year())
	}

	// month
	if r.Start.AddDate(0, 1, 0).Equal(*r.End) {
		y, m, _ := r.Start.Date()
		return fmt.Sprintf("%s %04d", r.locale.MonthWide(m), y)
	}

	// date
	if r.Start.AddDate(0, 0, 1).Equal(*r.End) {
		return r.locale.FmtDateShort(*r.Start)
	}

	// fallback: print the range, subtract one millisecond from end to be at the end of the day before (day a to day b does not include b)
	return fmt.Sprintf("%s - %s", r.locale.FmtDateShort(*r.Start), r.locale.FmtDateShort(r.End.Add(-1*time.Millisecond)))
}

func (r DateRange) Shift(years, months, days int) DateRange {
	var start, end time.Time
	if r.Start != nil && !r.Start.IsZero() {
		start = r.Start.AddDate(years, months, days)
	}
	if r.End != nil && !r.End.IsZero() {
		end = r.End.AddDate(years, months, days)
	}
	return NewDateRange(&start, &end, r.locale)
}

func (r DateRange) IsClosed() bool {
	return r.Start != nil && r.End != nil
}

func (r DateRange) IsMonthRange() bool {
	return r.End != nil && r.Start.AddDate(0, 1, 0) == *r.End
}

func (r DateRange) IsYearRange() bool {
	return r.End != nil && r.Start.AddDate(1, 0, 0) == *r.End
}

func (r DateRange) IsOpen() bool {
	return !r.IsClosed()
}

func (r DateRange) Empty() bool {
	return r.Start == nil && r.End == nil
}

func (r DateRange) Contains(date time.Time) bool {
	return r.Start != nil && !date.Before(*r.Start) && r.End != nil && !date.After(*r.End)
}

func (r DateRange) ContainsP(date *time.Time) bool {
	return date != nil && r.Contains(*date)
}

func (r DateRange) Intersection(start *time.Time, end *time.Time) time.Duration {
	if start == nil || end == nil || r.IsOpen() {
		return 0
	}

	containsStart := r.ContainsP(start)
	containsEnd := r.ContainsP(end)

	if containsStart && containsEnd {
		return end.Sub(*start)
	}

	if containsStart && !containsEnd {
		return r.End.Sub(*start)
	}

	if !containsStart && containsEnd {
		return end.Sub(*r.Start)
	}

	if start.Before(*r.Start) && end.After(*r.End) {
		return r.End.Sub(*r.Start)
	}

	return 0
}

func (r DateRange) Years(loc *time.Location) []DateRange {
	first := r.Start.In(loc).Year()
	last := r.End.In(loc).Year()

	var result []DateRange
	for i := first; i <= last; i++ {
		result = append(result, NewYearRange(time.Date(i, time.January, 1, 0, 0, 0, 0, loc), r.locale, loc))
	}
	return result
}

func (r DateRange) Months(loc *time.Location) []DateRange {
	end := r.End.In(loc)
	month := NewMonthRange(*r.Start, r.locale, loc)

	var result []DateRange
	for !month.Start.After(end) {
		result = append(result, month)
		month = month.Shift(0, 1, 0)
	}
	return result
}

func (r DateRange) Days(loc *time.Location) []DateRange {
	end := r.End.In(loc)
	month := NewDayRange(*r.Start, r.locale, loc)

	var result []DateRange
	for !month.Start.After(end) {
		result = append(result, month)
		month = month.Shift(0, 0, 1)
	}
	return result
}
