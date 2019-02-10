package property

import (
	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/util"
)

func newRemoveCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	force := false

	var cmd = &cobra.Command{
		Use:   "remove <property id or name> ...",
		Short: "Removes properties",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, idOrName := range args {
				if err := ctx.StoreHelper.RemoveProperty(idOrName, force); err != nil {
					util.Fatal(err)
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "", force, "Remove property values from projects")

	parent.AddCommand(cmd)
	return cmd
}
