package context

import (
	"github.com/go-playground/locales"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/jansorg/tom/go-tom/i18n"
	"github.com/jansorg/tom/go-tom/query"
	"github.com/jansorg/tom/go-tom/store"
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
