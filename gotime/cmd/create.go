package cmd

import (
	"github.com/spf13/cobra"
)

func newCreateCommand(context *GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "create new content",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	newCreateProjectCommand(context, cmd)
	newCreateTagCommand(context, cmd)

	parent.AddCommand(cmd)
	return cmd
}
