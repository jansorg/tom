package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
)

func newResetCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	valisArgs := []string{"all", "projects", "tags", "frames"}

	var cmdReset = &cobra.Command{
		Use:       "reset [all | projects | tags | frames]",
		Short:     "resets the local database. Removes projects, tags and frames",
		ValidArgs: valisArgs,
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			arg := args[0]
			projects := arg == "all" || arg == "projects"
			tags := arg == "all" || arg == "tags"
			frames := arg == "all" || arg == "frames"

			if err := context.Store.Reset(projects, tags, frames); err != nil {
				fatal(err)
			}
		},
	}

	parent.AddCommand(cmdReset)
	return cmdReset
}
