package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
)

func newFramesCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	jsonOutput := false
	format := ""
	delimiter := ""

	var cmd = &cobra.Command{
		Use:   "frames",
		Short: "Print a listing of all frames",
		Run: func(cmd *cobra.Command, args []string) {
			frames := context.Store.Frames()

			if jsonOutput {
				if bytes, err := json.MarshalIndent(frames, "", "  "); err != nil {
					fatal(err)
				} else {
					fmt.Println(string(bytes))
				}
			} else {
				properties := strings.Split(format, ",")

				for _, frame := range frames {
					line := ""
					for i, prop := range properties {
						if i > 0 {
							line += delimiter
						}

						switch strings.TrimSpace(prop) {
						case "id":
							line += frame.ID
						case "projectID":
							line += frame.ProjectId
						case "projectName":
							project, _ := ctx.Query.ProjectByID(frame.ProjectId)
							line += project.FullName
						case "startTime":
							line += frame.Start.String()
						case "stopTime":
							line += frame.Start.String()
						case "duration":
							if frame.IsStopped() {
								line += ctx.DurationPrinter.Short(frame.Duration())
							} else {
								line += ""
							}
						default:
							fatal("unknown property", prop, ". Valid values: id, projectID, projectName, startTime, stopTime, duration")
						}
					}
					fmt.Println(line)
				}
			}
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "id", "A comma separated list of of properties to output. Default: id . Possible values: id,name,shortName")
	cmd.Flags().StringVarP(&delimiter, "delimiter", "d", "\t", "The delimiter to add between property values. Default: TAB")
	cmd.Flags().BoolVarP(&jsonOutput, "json", "", false, "Prints JSON instead of plain text")

	parent.AddCommand(cmd)
	return cmd
}
