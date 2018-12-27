package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var cmdProjects = &cobra.Command{
	Use:   "projects",
	Short: "Print a listing of all projects",
	Run: func(cmd *cobra.Command, args []string) {
		projects := Store.Projects()

		if jsonOutput {
			if bytes, err := json.Marshal(projects); err != nil {
				fatal(err)
			} else {
				fmt.Println(string(bytes))
			}
		} else {
			for _, p := range projects {
				fmt.Println(p.Id)
			}
		}
	},
}
