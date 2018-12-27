package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/store"
)

var cfgFile string
var dataDir string

var ctx context.GoTimeContext

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "data directory (default is $HOME/.gotime)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotime.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&ctx.JsonOutput, "json", "j", false, "output JSON instead of plain text")

	newProjectsCommand(&ctx, rootCmd)
	newFramesCommand(&ctx, rootCmd)
	newCreateCommand(&ctx, rootCmd)
	newRemoveCommand(&ctx, rootCmd)
	newStartCommand(&ctx, rootCmd)
	newStopCommand(&ctx, rootCmd)
	newReportCommand(&ctx, rootCmd)
	newResetCommand(&ctx, rootCmd)

	viper.BindPFlag("data-dir", rootCmd.PersistentFlags().Lookup("data-dir"))
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
}

var rootCmd = &cobra.Command{
	Use:   "gotime",
	Short: "gotime is a command line application to track time.",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
