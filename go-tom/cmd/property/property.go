package property

import (
	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
)

func NewCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "property",
		Short: "Manage the list of properties and the values applied to projects",
	}

	newCreateCommand(ctx, cmd)
	newRemoveCommand(ctx, cmd)
	newGetValueCommand(ctx, cmd)
	newSetValueCommand(ctx, cmd)

	parent.AddCommand(cmd)
	return cmd
}
