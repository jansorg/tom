package cmd

import (
	"github.com/spf13/cobra"

	"../context"
)

func newCreateCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "create new content",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	newCreateProjectCommand(ctx, cmd)
	newCreateTagCommand(ctx, cmd)

	parent.AddCommand(cmd)
	return cmd
}
