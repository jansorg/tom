package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/store"
)

func newRemoveProjectCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "project <project name or project ID> ...",
		Short: "removes new project and all its associated data",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			remvoedFrames := 0
			remvoedProjects := 0

			for _, name := range args {
				projects := context.Store.FindProjects(func(p store.Project) bool {
					return p.Id == name || p.FullName == name
				})

				for _, p := range projects {
					frames := context.Store.FindFrames(func(frame store.Frame) bool {
						return frame.ProjectId == p.Id
					})

					for _, f := range frames {
						if err := context.Store.RemoveFrame(f.Id); err != nil {
							fatal(err)
						}
						remvoedFrames++
					}

					if err := context.Store.RemoveProject(p.Id); err != nil {
						fatal(err)
					}
					remvoedProjects++
				}
			}

			fmt.Printf("Removed projects: %d\nRemoved frames: %d\n", remvoedProjects, remvoedFrames)
		},
	}
	parent.AddCommand(cmd)

	return cmd
}
