package property

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/properties"
	"github.com/jansorg/tom/go-tom/util"
)

func newCreateCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	typeName := properties.CurrencyType.ID()
	description := ""
	applyToSubprojects := true

	var cmd = &cobra.Command{
		Use:   "create <property name>",
		Short: "Adds a new property definition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			propType := properties.FindType(typeName)
			if propType == nil {
				util.Fatal("property type %s not found", typeName)
			}

			prop, err := ctx.Store.AddProperty(&properties.Property{
				Name:               args[0],
				Description:        description,
				TypeID:             propType.ID(),
				ApplyToSubprojects: applyToSubprojects,
			})

			if err != nil {
				util.Fatal(err)
			} else {
				fmt.Printf("created property %s with ID %s\n", prop.Name, prop.ID)
			}
		},
	}

	cmd.Flags().StringVarP(&typeName, "type", "", typeName, "Property data type. Values: currency")
	cmd.Flags().StringVarP(&description, "description", "", description, "Optional description with details about this property.")
	cmd.Flags().BoolVarP(&applyToSubprojects, "subprojects", "s", applyToSubprojects, "Values of this property will be inherited by subprojects")

	parent.AddCommand(cmd)
	return cmd
}
