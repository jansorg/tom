package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/util"
)

func newProjectsPropertyCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "property",
		Short: "Get/set project properties. Usage: <projectName> [propertyName] --set [optional new value]",
		Args:  cobra.RangeArgs(1, 3),
		Run: func(cmd *cobra.Command, args []string) {
			project, err := ctx.Query.ProjectByFullName(strings.Split(args[0], "/"))
			if err != nil {
				util.Fatal(err)
			}

			if len(args) == 1 {
				size := len(project.Properties)
				fmt.Printf("%d properties for %s\n\n", size, args[0])
				if size > 0 {
					fmt.Println("Properties:")
					for k, v := range project.Properties {
						fmt.Printf("\t%s=%v\n", k, v)
					}
				}

				inherited := 0
				out := ""
				ctx.Query.WithProjectAndParents(project.ID, func(parent *model.Project) bool {
					if project.ID != parent.ID {
						inherited += len(parent.Properties)

						for k, v := range parent.Properties {
							out += fmt.Sprintf("\t%s=%v (from %s)\n", k, v, parent.FullName)
						}
					}
					return true
				})
				if inherited > 0 {
					fmt.Println("Inherited properties:")
					fmt.Println(out)
				}
			} else if len(args) == 2 {
				fmt.Printf("%s=%v\n", args[1], project.Properties[args[1]])
			} else if len(args) == 3 {
				old := project.Properties[args[1]]
				project.Properties[args[1]] = args[2]
				ctx.Store.UpdateProject(*project)
				fmt.Printf("%s=%v (previously: %s)\n", args[1], args[2], old)
			} else {
				util.Fatal(fmt.Errorf("unsupported configuration"))
			}
		},
	}

	parent.AddCommand(cmd)
	return cmd
}
