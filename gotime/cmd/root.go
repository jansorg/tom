package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/message"

	"github.com/jansorg/gotime/gotime/config"
	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/i18n"
	"github.com/jansorg/gotime/gotime/query"
	"github.com/jansorg/gotime/gotime/store"
)

var ctx context.GoTimeContext
var configFile string

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().String("data-dir", "", "data directory (default is $HOME/.gotime)")
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is $HOME/.gotime.yaml)")

	newProjectsCommand(&ctx, RootCmd)
	newFramesCommand(&ctx, RootCmd)
	newCreateCommand(&ctx, RootCmd)
	newRemoveCommand(&ctx, RootCmd)
	newStartCommand(&ctx, RootCmd)
	newStopCommand(&ctx, RootCmd)
	newReportCommand(&ctx, RootCmd)
	newImportCommand(&ctx, RootCmd)
	newResetCommand(&ctx, RootCmd)
	newStatusCommand(&ctx, RootCmd)
	newCompletionCommand(&ctx, RootCmd)
	newConfigCommand(&ctx, RootCmd)

	if err := viper.BindPFlag(config.KeyDataDir, RootCmd.PersistentFlags().Lookup("data-dir")); err != nil {
		fatal(err)
	}
}

func fatal(err ...interface{}) {
	fmt.Println(append([]interface{}{"Error: "}, err...)...)
	os.Exit(1)
}

func initConfig() {
	config.SetDefaults()
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	// setup config dir if it doesn't exist
	dataDir := viper.GetString(config.KeyDataDir)
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0700); err != nil {
			fatal(err)
		}
	}

	if err := viper.ReadInConfig(); !os.IsNotExist(err) {
		// fatal(err)
	}

	dataStore, err := store.NewStore(dataDir)
	if err != nil {
		fatal(err)
	}

	ctx.Store = dataStore
	ctx.StoreHelper = store.NewStoreHelper(dataStore)
	ctx.Query = query.NewStoreQuery(dataStore)
	ctx.Language = i18n.FindPreferredLanguages()
	ctx.LocalePrinter = message.NewPrinter(ctx.Language)
	ctx.Locale = i18n.FindLocale(ctx.Language)
	ctx.DurationPrinter = i18n.NewDurationPrinter(ctx.Language)
	ctx.DateTimePrinter = i18n.NewDateTimePrinter(ctx.Language)
}

var RootCmd = &cobra.Command{
	Use:     "gotime",
	Short:   "gotime is a command line application to track time.",
	Version: "1.0.0",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
