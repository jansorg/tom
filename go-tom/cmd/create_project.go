package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/go-tom/context"
)

func newCreateProjectCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "project",
		Short:   "Create a new project",
		Args:    cobra.MinimumNArgs(1),
		Example: "gotime create project \"Installation\" \"Deployment\"",
		Run: func(cmd *cobra.Command, args []string) {
			for _, name := range args {
				project, created, err := context.StoreHelper.GetOrCreateNestedProject(name)
				if err != nil {
					fatal(err)
				}

				if created {
					fmt.Printf("created project %s\n", project.FullName)
				} else {
					fmt.Printf("using existing project %s\n", project.FullName)
				}
			}
		},
	}
	parent.AddCommand(cmd)

	return cmd
}
