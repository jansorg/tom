package property

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/util"
)

func newCreateCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	typeName := "number"
	applyToSubprojects := true
	prefix := ""
	suffix := ""

	var cmd = &cobra.Command{
		Use:   "create <property name>",
		Short: "Adds a new property definition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			typeID, err := model.TypeFromString(typeName)
			if err != nil {
				util.Fatal(err)
			}

			prop, err := ctx.Store.AddProperty(&model.Property{
				Name:               args[0],
				Type:               typeID,
				Prefix:             prefix,
				Suffix:             suffix,
				ApplyToSubprojects: applyToSubprojects,
			})

			if err != nil {
				util.Fatal(err)
			} else {
				fmt.Printf("created property %s with ID %s\n", prop.Name, prop.ID)
			}
		},
	}

	cmd.Flags().StringVarP(&typeName, "type", "", typeName, "Property data type. Values: string | number")
	cmd.Flags().BoolVarP(&applyToSubprojects, "subprojects", "s", applyToSubprojects, "Values of this property will be inherited by subprojects")
	cmd.Flags().StringVarP(&prefix, "prefix", "", prefix, "Prefix string to display next to values of this property type")
	cmd.Flags().StringVarP(&suffix, "suffix", "", suffix, "Suffix string to display next to values of this property type")

	parent.AddCommand(cmd)
	return cmd
}
