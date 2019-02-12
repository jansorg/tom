package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jansorg/tom/go-tom/config"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/util"
)

func newConfigSetCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "set",
		Short: "Sets a configuration value",
		Args:  cobra.ExactArgs(2),
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
