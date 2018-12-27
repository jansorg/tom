package cmd

import (
	"github.com/spf13/cobra"
)

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "create new content",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cmdCreate.AddCommand(cmdCreateProject)
}
