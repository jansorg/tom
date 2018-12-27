package cmd

import (
	"github.com/spf13/cobra"

	"../store"
)

var cmdCreateProject = &cobra.Command{
	Use:   "project",
	Short: "create a new project",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, name := range args {
			if err := Store.AddProject(store.Project{FullName: name, ShortName: name}); err != nil {
				fatal(err)
			}
		}
	},
}
