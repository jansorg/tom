package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"

	"github.com/jansorg/tom/go-tom/config"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/util"
)

func NewCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	output := "yaml"

	var cmd = &cobra.Command{
		Use:       "config",
		Short:     "Output configuration values as YAML or JSON.",
		Long:      "If no arguments are passed, then the complete configuration will be printed. If one or more arguments are passed, then each is printed with its current configuration values. Bash completion will suggest built-in configuration keys.",
		ValidArgs: config.Keys,
		Run: func(cmd *cobra.Command, args []string) {
			data, err := doConfigCommand(output, args...)
			if err != nil {
				util.Fatal(err)
			} else {
				fmt.Println(string(data))
			}
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", output, "Output format. Supported: yaml | json")

	newConfigSetCommand(ctx, cmd)
	parent.AddCommand(cmd)
	return cmd
}

func doConfigCommand(outputFormat string, keys ...string) ([]byte, error) {
	settings := make(map[string]interface{})
	if len(keys) == 0 {
		settings = viper.AllSettings()
	} else {
		for _, k := range keys {
			settings[k] = viper.Get(k)
		}
	}

	switch outputFormat {
	case "yaml":
		return yaml.Marshal(settings)
	case "json":
		return json.MarshalIndent(settings, "", "  ")
	default:
		return nil, fmt.Errorf("unsupported format %s", outputFormat)
	}
}

func createConfigIfNotExists(dataDir string) error {
	filePath := filepath.Join(dataDir, fmt.Sprintf("%s.yaml", config.ConfigFilename))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ioutil.WriteFile(filePath, []byte{}, 0600)
	}
	return nil
}
