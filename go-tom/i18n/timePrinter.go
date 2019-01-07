package i18n

import (
	"time"

	"golang.org/x/text/language"
)

type DateTimePrinter interface {
	Time(d time.Time) string
	TimeShort(d time.Time) string

	Date(d time.Time) string

	DateTime(d time.Time) string
	DateTimeShort(d time.Time) string
}

func NewDateTimePrinter(lang language.Tag) DateTimePrinter {
	return &dateTimePrinter{lang: lang}
}

type dateTimePrinter struct {
	lang language.Tag
}

func (p *dateTimePrinter) Time(d time.Time) string {
	return d.Format("15:04:05")
}

func (p *dateTimePrinter) TimeShort(d time.Time) string {
	return d.Format("15:04")
}

func (p *dateTimePrinter) Date(d time.Time) string {
	return d.Format("2006-01-02")
}

func (p *dateTimePrinter) DateTime(d time.Time) string {
	return d.Format("2006-01-02 15:04:05")
}

func (p *dateTimePrinter) DateTimeShort(d time.Time) string {
	return d.Format("2006-01-02 15:04")
}
