package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
)

func newCreateInvoiceCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "invoice",
		Short:   "Create a new invoice, based on the tracked time of a project",
		Args:    cobra.MinimumNArgs(1),
		Example: "gotime create invoice",
		Run: func(cmd *cobra.Command, args []string) {
			for _, name := range args {
				project, err := context.StoreHelper.GetOrCreateNestedProject(name)
				if err != nil {
					fatal(err)
				}

				fmt.Printf("created invoice %s\n", project.FullName)
			}
		},
	}
	parent.AddCommand(cmd)

	return cmd
}
