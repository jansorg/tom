package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
)

func newRemoveAllCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	validArgs := []string{"all", "projects", "tags", "frames"}

	var cmd = &cobra.Command{
		Use:       "all [all | projects | tags | frames]",
		Short:     "Removes all stores data. Specify the type to only remove projects, tags or frames",
		ValidArgs: validArgs,
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			arg := args[0]
			projects := arg == "all" || arg == "projects"
			tags := arg == "all" || arg == "tags"
			frames := arg == "all" || arg == "frames"

			if removedProjects, removedTags, removedFrames, err := context.Store.Reset(projects, tags, frames); err != nil {
				fatal(err)
			} else {
				fmt.Printf("Successfully removed %d projects, %d tags and %d frames\n", removedProjects, removedTags, removedFrames)
			}
		},
	}
	parent.AddCommand(cmd)
	return cmd
}
