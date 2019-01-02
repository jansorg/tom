package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
)

func newImportCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "import",
		Short: "import frames and projects",
	}

	newImportFanurioCommand(ctx, cmd)
	newImportWatsonCommand(ctx, cmd)

	parent.AddCommand(cmd)
	return cmd
}
