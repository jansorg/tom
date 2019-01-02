package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
)

func newRemoveCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "remove [all projects | tags | frames] or remove project <project name>",
		Short: "remove projects, tags or frames",
	}

	newRemoveProjectCommand(ctx, cmd)
	newRemoveAllCommand(ctx, cmd)

	parent.AddCommand(cmd)
	return cmd
}
