package edit

import (
	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
)

func NewEditCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "edit project | frame",
		Short: "edit properties of projects or frames",
	}

	newEditFrameCommand(ctx, cmd)
	newEditProjectCommand(ctx, cmd)

	parent.AddCommand(cmd)
	return cmd
}
