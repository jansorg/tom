package remove

import (
	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
)

func NewCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "remove [all projects | tags | frames] or remove project <project name> or frame ID",
		Short: "remove projects, tags or frames",
	}

	newRemoveProjectCommand(ctx, cmd)
	newRemoveFrameCommand(ctx, cmd)
	newRemoveAllCommand(ctx, cmd)

	parent.AddCommand(cmd)
	return cmd
}
