package context

import (
	"github.com/go-playground/locales"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/jansorg/gotime/gotime/query"
	"github.com/jansorg/gotime/gotime/store"
)

type GoTimeContext struct {
	Store        store.Store
	StoreHelper  *store.StoreHelper
	Query        query.StoreQuery
	JsonOutput   bool
	Language     language.Tag
	Translator   *message.Printer
	NumberFormat *message.Printer
	Locale       locales.Translator
}
