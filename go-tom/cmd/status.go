package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
)

func newStatusCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	verbose := false
	format := ""
	delimiter := ""

	var cmd = &cobra.Command{
		Use:   "status",
		Short: "Displays when the current project was started and the time spent...",
		Run: func(cmd *cobra.Command, args []string) {
			activeFrames := ctx.Query.ActiveFrames

			if cmd.Flag("format").Changed {
				flags := strings.Split(format, ",")

				for _, frame := range activeFrames() {
					project, err := ctx.Query.ProjectByID(frame.ProjectId)
					if err != nil {
						fatal(err)
					}

					var values []string
					for _, flag := range flags {
						value := ""
						switch flag {
						case "id":
							value = frame.ID
						case "projectID":
							value = project.ID
						case "projectName":
							value = project.Name
						case "projectFullName":
							value = project.FullName
						case "projectParentID":
							value = project.ParentID
						case "startTime":
							value = frame.Start.Format(time.RFC3339)
						default:
							fatal(fmt.Errorf("unknown flag %s", flag))
						}

						values = append(values, value)
					}
					fmt.Println(strings.Join(values, delimiter))
				}
			} else if verbose {
				projectCount := len(ctx.Store.Projects())
				tagCount := len(ctx.Store.Tags())
				frameCount := len(ctx.Store.Frames())
				activeFrameCount := len(activeFrames())

				fmt.Printf("Projects: %d\nTags: %d\nFrames: %d\nStarted activites: %d\n", projectCount, tagCount, frameCount, activeFrameCount)
			} else {
				for _, frame := range activeFrames() {
					project, err := ctx.Query.ProjectByID(frame.ProjectId)
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

	cmd.Flags().StringVarP(&format, "format", "f", "", "Properties to print for each active frame. Possible values: id,projectID,projectName,projectFullName,projectParentID,startTime")
	cmd.Flags().StringVarP(&delimiter, "delimiter", "d", "\t", "Delimiter to separate flags on the same line. Only used when --format is specified.")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Print details about the currently stored projects, tags and frames")

	newProjectsStatusCommand(ctx, cmd)
	parent.AddCommand(cmd)
	return cmd
}
