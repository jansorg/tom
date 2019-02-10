package property

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/cmd/cmdUtil"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/properties"
	"github.com/jansorg/tom/go-tom/util"
)

type propertyList struct {
	project    *model.Project
	properties []*properties.Property
}

func (p propertyList) Size() int {
	return len(p.properties)
}

func (p propertyList) Get(index int, propName string, format string, ctx *context.TomContext) (interface{}, error) {
	property := p.properties[index]

	switch propName {
	case "id":
		return property.ID, nil
	case "name":
		return property.Name, nil
	case "value":
		return ctx.Query.FindPropertyValue(property.ID, p.project.ID)
	default:
		return nil, fmt.Errorf("unknown property %s", propName)
	}
}

func newGetValueCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	nameDelimiter := ""

	var cmd = &cobra.Command{
		Use:   "get <project id or name> [property id or name]",
		Short: "Get the value of project property.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			project, err := ctx.Query.ProjectByFullNameOrID(args[0], nameDelimiter)
			if err != nil {
				util.Fatal(err)
			}

			data := propertyList{project: project}

			switch len(args) {
			case 1:
				values := ctx.Query.FindPropertyValues(project.ID)
				for prop := range values {
					data.properties = append(data.properties, prop)
				}
			case 2:
				if property, err := ctx.Query.FindPropertyByNameOrID(args[1]); err != nil {
					util.Fatal(err)
				} else {
					data.properties = []*properties.Property{property}
				}
			default:
				util.Fatal("unsupported number of arguments")
			}

			if err := cmdUtil.PrintList(cmd, data, ctx); err != nil {
				util.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&nameDelimiter, "name-delimiter", "", "/", "Delimiter used in the full project name")
	cmdUtil.AddListOutputFlags(cmd, "id,name,value", []string{"id", "name", "value"})

	parent.AddCommand(cmd)
	return cmd
}
