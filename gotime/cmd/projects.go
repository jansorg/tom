package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
)

func newProjectsCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	format := ""
	delimiter := ""
	jsonOutput := false

	var cmd = &cobra.Command{
		Use:   "projects",
		Short: "Print a listing of all projects",
		Run: func(cmd *cobra.Command, args []string) {
			projects := context.Store.Projects()
			sort.SliceStable(projects, func(i, j int) bool {
				return strings.Compare(projects[i].FullName, projects[j].FullName) < 0
			})

			properties := strings.Split(format, ",")

			if jsonOutput {
				if bytes, err := json.MarshalIndent(projects, "", "  "); err != nil {
					fatal(err)
				} else {
					fmt.Println(string(bytes))
				}
			} else {
				for _, p := range projects {
					line := ""
					for i, prop := range properties {
						if i > 0 {
							line += delimiter
						}
						switch strings.TrimSpace(prop) {
						case "id":
							line += p.ID
						case "name":
							line += p.FullName
						case "shortName":
							line += p.Name
						default:
							fatal("unknown property", prop)
						}
					}
					fmt.Println(line)
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&jsonOutput, "json", "", false, "Prints JSON instead of plain text")
	cmd.Flags().StringVarP(&format, "format", "f", "id", "A comma separated list of of properties to output. Default: id . Possible values: id,name,shortName")
	cmd.Flags().StringVarP(&delimiter, "delimiter", "d", "\t", "The delimiter to add between property values. Default: TAB")

	parent.AddCommand(cmd)
	return cmd
}
