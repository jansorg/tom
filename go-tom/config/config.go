package config

import (
	"log"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const KeyDataDir = "data_dir"
const KeyActivityStopOnStart = "activity.stop_on_start"
const KeyProjectCreateMissing = "projects.create_missing"

const ConfigFilename = "gotime"

func SetDefaults() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	dataDirPath := filepath.Join(home, ".gotime")
	viper.SetDefault(KeyDataDir, dataDirPath)
	viper.SetDefault(KeyProjectCreateMissing, false)
	viper.SetDefault(KeyActivityStopOnStart, true)

	viper.SetConfigName(ConfigFilename)
	// fixme add /etc?
	viper.AddConfigPath(dataDirPath)
}
