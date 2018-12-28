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
	showNames := false

	var cmd = &cobra.Command{
		Use:   "projects",
		Short: "Print a listing of all projects",
		Run: func(cmd *cobra.Command, args []string) {
			projects := context.Store.Projects()
			sort.SliceStable(projects, func(i, j int) bool {
				return strings.Compare(projects[i].FullName, projects[j].FullName) < 0
			})

			if context.JsonOutput {
				if bytes, err := json.MarshalIndent(projects, "", "  "); err != nil {
					fatal(err)
				} else {
					fmt.Println(string(bytes))
				}
			} else {
				for _, p := range projects {
					if showNames {
						fmt.Printf("%s\t%s\n", p.ID, p.FullName)
					} else {
						fmt.Println(p.ID)
					}
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&showNames, "names", "n", false, "Display the names in the plain text output")

	parent.AddCommand(cmd)
	return cmd
}
