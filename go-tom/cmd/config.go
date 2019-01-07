package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"

	"github.com/jansorg/gotime/go-tom/config"
	"github.com/jansorg/gotime/go-tom/context"
)

func newConfigCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "config",
		Short: "prints the current configuration on stdout",
		Long:  "Prints the configuration values of gotime. If no arguments are passed, then the complete configuration will be printed. If one or more arguments are passed, then each is printed with its current configuration values.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				if cfg, err := settings(); err != nil {
					fatal(err)
				} else {
					fmt.Println(cfg)
				}
			} else {
				for _, name := range args {
					value := viper.Get(name)
					bs, err := yaml.Marshal(value)
					if err != nil {
						fatal(fmt.Errorf("unable to marshal config to YAML: %v", err))
					}
					fmt.Printf("%s=%s", name, string(bs))
				}
			}
		},
	}

	newConfigSetCommand(ctx, cmd)
	parent.AddCommand(cmd)

	return cmd
}

func settings() (string, error) {
	c := viper.AllSettings()

	bs, err := yaml.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("unable to marshal config to YAML: %v", err)
	}
	return string(bs), nil
}

func createEmptyConfig(dataDir string) error {
	filePath := filepath.Join(dataDir, fmt.Sprintf("%s.yaml", config.ConfigFilename))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ioutil.WriteFile(filePath, []byte{}, 0600)
	}
	return nil
}
