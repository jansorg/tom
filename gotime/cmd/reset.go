package cmd

import (
	"github.com/spf13/cobra"
)

var cmdReset = &cobra.Command{
	Use:   "reset",
	Short: "resets the local database. Removes projects, tags and frames",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Store.Reset(); err != nil {
			fatal(err)
		}
	},
}
