package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/cmd/util"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

func newCreateProjectCommand(context *context.TomContext, parent *cobra.Command) *cobra.Command {
	var parentID string

	var output string

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
					util.Fatal(err)
				}

				if !created {
					util.Fatal(fmt.Printf("the project %s does already exist\n", project.FullName))

				} else if output == "json" {
					util.PrintJSON((*model.DetailedProject)(project))
				} else {
					fmt.Printf("created project %s\n", project.FullName)
				}
			}
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "plain", "Output format. Supported: plain | json. Default: plain")

	cmd.Flags().StringVarP(&parentID, "parent", "p", "", "Optional parent project ID. If defined the new project will be made a child project of this. If defined the name will be used as is for the name of the new project.")
	parent.AddCommand(cmd)

	return cmd
}
