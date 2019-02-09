package property

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/cmd/cmdUtil"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/util"
)

func newSetValueCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	nameDelimiter := ""

	var cmd = &cobra.Command{
		Use:   "set <project id or name> <property name or id> <value>",
		Short: "Set a property value for a project.",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			project, err := ctx.Query.ProjectByFullNameOrID(args[0], nameDelimiter)
			if err != nil {
				util.Fatal(err)
			}

			property, err := ctx.Query.FindPropertyByNameOrID(args[1])
			if err != nil {
				util.Fatal(err)
			}

			previousValue, _ := ctx.Query.FindPropertyValue(property.ID, project.ID)
			convertedValue, err := property.FromString(args[2])
			if err != nil {
				util.Fatal(err)
			}

			if err = project.SetPropertyValue(property.ID, convertedValue); err != nil {
				util.Fatal(err)
			}

			if _, err := ctx.Store.UpdateProject(*project); err != nil {
				util.Fatal(err)
			}

			if previousValue == nil {
				fmt.Printf("%s=%v\n", args[1], convertedValue)
			} else {
				fmt.Printf("%s=%v (previously: %s)\n", args[1], convertedValue, previousValue)
			}
		},
	}

	cmd.Flags().StringVarP(&nameDelimiter, "name-delimiter", "", "/", "Delimiter used in the full project name")
	cmdUtil.AddListOutputFlags(cmd, "id,name,value", []string{"id", "name", "value"})

	parent.AddCommand(cmd)
	return cmd
}
