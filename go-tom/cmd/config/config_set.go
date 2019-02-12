package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jansorg/tom/go-tom/config"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/util"
)

func newSetCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "set <key> <value>",
		Short:                 "Sets a configuration value",
		Long:                  "Sets a configuration value to a new value. Use `tom config` to see a list of all supported keys. Bash completion will suggest the names of built-in keys.",
		Example:               fmt.Sprintf("tom config set %s $HOME/tom-data", config.KeyDataDir),
		Args:                  cobra.ExactArgs(2),
		ValidArgs:             config.Keys,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			if err := createConfigIfNotExists(viper.GetString(config.KeyDataDir)); err != nil {
				util.Fatal(err)
			}

			viper.Set(args[0], args[1])
			if err := viper.WriteConfig(); err != nil {
				util.Fatal("error updating configuration file: ", err)
			}

			fmt.Println("Successfully updated the configuration value of " + args[0])
		},
	}

	parent.AddCommand(cmd)
	return cmd
}
