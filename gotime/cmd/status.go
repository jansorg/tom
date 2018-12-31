package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
)

func newStatusCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	verbose := false

	var cmd = &cobra.Command{
		Use:   "status",
		Short: "Displays when the current project was started and the time spent...",
		Run: func(cmd *cobra.Command, args []string) {
			activeFrames := context.Query.ActiveFrames

			if verbose {
				projectCount := len(context.Store.Projects())
				tagCount := len(context.Store.Tags())
				frameCount := len(context.Store.Frames())
				activeFrameCount := len(activeFrames())

				fmt.Printf("Projects: %d\nTags: %d\nFrames: %d\nStarted activites: %d\n", projectCount, tagCount, frameCount, activeFrameCount)
			} else {
				for _, frame := range activeFrames() {
					project, err := context.Query.ProjectByID(frame.ProjectId)
					if err != nil {
						fatal(err)
					}

					fmt.Printf("Project %s was started %s\n", project.FullName, ctx.DateTimePrinter.DateTime(*frame.Start))
				}

				if len(activeFrames()) == 0 {
					fmt.Printf("%d active frames found\n", len(activeFrames()))
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Print details about the currently stored projects, tags and frames")

	parent.AddCommand(cmd)
	return cmd
}
