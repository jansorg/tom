package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/store"
)

func newCreateProjectCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "project",
		Short: "create a new project",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, name := range args {
				if _, err := context.Store.AddProject(store.Project{FullName: name, ShortName: name}); err != nil {
					fatal(err)
				}
			}
		},
	}
	parent.AddCommand(cmd)

	return cmd
}
