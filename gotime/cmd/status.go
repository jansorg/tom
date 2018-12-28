package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
)

func newStatusCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "status",
		Short: "prints the current status",
		Run: func(cmd *cobra.Command, args []string) {
			projectCount := len(context.Store.Projects())
			tagCount := len(context.Store.Tags())
			frameCount := len(context.Store.Frames())

			fmt.Printf("Projects: %d\nTags: %d\nFrames: %d\n", projectCount, tagCount, frameCount)
		},
	}

	parent.AddCommand(cmd)
	return cmd
}
