package test_setup

import (
	"io/ioutil"
	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/i18n"
	"github.com/jansorg/tom/go-tom/query"
	store2 "github.com/jansorg/tom/go-tom/store"
	"github.com/jansorg/tom/go-tom/storeHelper"
)

func CreateTestContext(lang language.Tag) (*context.TomContext, error) {
	dir, err := ioutil.TempDir("", "tom-data")
	if err != nil {
		return nil, err
	}

	backupDir, err := ioutil.TempDir("", "tom-backup")
	if err != nil {
		return nil, err
	}

	store, err := store2.NewStore(dir, backupDir, 5)
	if err != nil {
		return nil, err
	}

	return &context.TomContext{
		Store:           store,
		StoreHelper:     storeHelper.NewStoreHelper(store),
		Query:           query.NewStoreQuery(store),
		Language:        lang,
		LocalePrinter:   message.NewPrinter(lang),
		Locale:          i18n.FindLocale(lang, false),
		DurationPrinter: i18n.NewDurationPrinter(lang),
	}, nil
}

func CleanupTestContext(ctx *context.TomContext) {
	os.RemoveAll(ctx.Store.DirPath())
	os.RemoveAll(ctx.Store.BackupDirPath())
}
