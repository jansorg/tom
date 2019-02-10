package imports

import (
	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
)

func NewCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "import",
		Short: "import frames and projects",
	}

	newFanurioCommand(ctx, cmd)
	newWatsonCommand(ctx, cmd)
	newMacTimeTrackCommand(ctx, cmd)

	parent.AddCommand(cmd)
	return cmd
}
