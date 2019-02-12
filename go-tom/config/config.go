package config

import (
	"log"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const KeyDataDir = "data_dir"
const KeyBackupDir = "backup_dir"
const KeyMaxBackups = "max_backups"
const KeyActivityStopOnStart = "activity.stop_on_start"
const KeyProjectCreateMissing = "projects.create_missing"

const ConfigFilename = "tom"

func SetDefaults() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	dataDirPath := filepath.Join(home, ".tom")
	backupDirPath := filepath.Join(dataDirPath, "backup")
	viper.SetDefault(KeyDataDir, dataDirPath)
	viper.SetDefault(KeyBackupDir, backupDirPath)
	viper.SetDefault(KeyMaxBackups, 10)
	viper.SetDefault(KeyProjectCreateMissing, false)
	viper.SetDefault(KeyActivityStopOnStart, true)

	viper.SetConfigName(ConfigFilename)
	// fixme add /etc?
	viper.AddConfigPath(dataDirPath)
}
