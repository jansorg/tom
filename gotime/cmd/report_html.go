package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
)

func newReportHtmlCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var templateFile string

	var cmd = &cobra.Command{
		Use:   "html",
		Short: "Generate a HTML report",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmd.Flags().StringVarP(&templateFile, "template", "t", "default", "The Go template to use. Default: templates/default.gohtml")

	parent.AddCommand(cmd)
	return cmd
}
