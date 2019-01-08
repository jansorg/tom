package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

func newCreateTagCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "tag name...",
		Short: "create a new tag",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, name := range args {
				if _, err := context.Store.AddTag(model.Tag{Name: name}); err != nil {
					fatal(err)
				}
			}
		},
	}
	parent.AddCommand(cmd)

	return cmd
}
