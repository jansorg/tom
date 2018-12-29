package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/message"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/i18n"
	"github.com/jansorg/gotime/gotime/query"
	"github.com/jansorg/gotime/gotime/store"
)

var cfgFile string
var dataDir string

var ctx context.GoTimeContext

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "data directory (default is $HOME/.gotime)")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotime.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&ctx.JsonOutput, "json", "j", false, "output JSON instead of plain text")

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

	viper.BindPFlag("data-dir", RootCmd.PersistentFlags().Lookup("data-dir"))
}

func fatal(err ...interface{}) {
	fmt.Println(err...)
	os.Exit(1)
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(home)
		viper.SetConfigName(".gotime")
	}

	viper.SetDefault("data-dir", filepath.Join(home, ".gotime"))

	if err := viper.ReadInConfig(); err != nil {
		// fmt.Println("Can't read config:", err)
		// os.Exit(1)
	}

	dataDir := viper.GetString("data-dir")
	os.MkdirAll(dataDir, 0700)
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
