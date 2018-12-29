package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
)

func newCreateCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "create new content. See the available sub commands.",
	}

	newCreateProjectCommand(ctx, cmd)
	newCreateTagCommand(ctx, cmd)

	parent.AddCommand(cmd)
	return cmd
}
