package context

import (
	"github.com/go-playground/locales"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/jansorg/gotime/gotime/i18n"
	"github.com/jansorg/gotime/gotime/query"
	"github.com/jansorg/gotime/gotime/store"
)

type GoTimeContext struct {
	Store           store.Store
	StoreHelper     *store.Helper
	Query           query.StoreQuery
	Language        language.Tag
	LocalePrinter   *message.Printer
	Locale          locales.Translator
	DurationPrinter i18n.DurationPrinter
	DateTimePrinter i18n.DateTimePrinter
}
