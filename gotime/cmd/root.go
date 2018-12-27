package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"../store"
)

var cfgFile string
var dataDir string

type GoTimeContext struct {
	Store      store.Store
	JsonOutput bool
}

var context GoTimeContext

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "data directory (default is $HOME/.gotime)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotime.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&context.JsonOutput, "json", "j", false, "output JSON instead of plain text")

	newProjectsCommand(&context, rootCmd)
	newFramesCommand(&context, rootCmd)
	newCreateCommand(&context, rootCmd)
	newStartCommand(&context, rootCmd)
	newStopCommand(&context, rootCmd)
	newResetCommand(&context, rootCmd)

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

	context.Store = dataStore
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
