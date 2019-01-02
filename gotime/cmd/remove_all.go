package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
)

func newRemoveAllCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	validArgs := []string{"all", "projects", "tags", "frames"}

	var cmd = &cobra.Command{
		Use:       "all [projects | tags | frames]",
		Short:     "Removes all stores data. Specify the type to only remove projects, tags or frames",
		ValidArgs: validArgs,
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			arg := args[0]
			projects := arg == "projects"
			tags := arg == "tags"
			frames := arg == "frames"

			if err := context.Store.Reset(projects, tags, frames); err != nil {
				fatal(err)
			}

			fmt.Printf("Successfully removed all %s\n", arg)
		},
	}
	parent.AddCommand(cmd)
	return cmd
}
