package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

func newCreateProjectCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var parentID string

	var cmd = &cobra.Command{
		Use:     "project",
		Short:   "Create a new project",
		Args:    cobra.MinimumNArgs(1),
		Example: "gotime create project \"Installation\" \"Deployment\" \"Support\"",
		Run: func(cmd *cobra.Command, args []string) {
			for _, name := range args {
				var project *model.Project
				var created bool
				var err error

				if parentID != "" {
					project, created, err = context.StoreHelper.GetOrCreateProject(name, parentID)
				} else {
					project, created, err = context.StoreHelper.GetOrCreateNestedProject(name)
				}

				if err != nil {
					Fatal(err)
				}

				if created {
					fmt.Printf("created project %s\n", project.FullName)
				} else {
					fmt.Printf("the project %s does already exist\n", project.FullName)
					// os.Exit(1)
				}
			}
		},
	}

	cmd.Flags().StringVarP(&parentID, "parent", "p", "", "Optional parent project ID. If defined the new project will be made a child project of this. If defined the name will be used as is for the name of the new project.")
	parent.AddCommand(cmd)

	return cmd
}
