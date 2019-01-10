package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
)

func newRenameCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rename TYPE oldName newName",
		Short: "rename a project or a tag",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var typeName = args[0]
			var oldName = args[1]
			var newName = args[2]

			switch typeName {
			case "project":
				if _, err := ctx.StoreHelper.RenameProjectByName(oldName, newName); err != nil {
					fatalf("rename failed: %s", err.Error())
				}
			case "tag":
				if tag, err := ctx.Query.TagByName(oldName); err != nil {
					fatal("tag %s not found", oldName)
				} else {
					tag.Name = newName
					if _, err := ctx.Store.UpdateTag(*tag); err != nil {
						fatal("unable to rename tag %s to %s", oldName, newName)
					}
				}
			default:
				fatal(fmt.Errorf("unknown TYPE %s. Valid values are project, tag", typeName))
			}

			fmt.Printf("successfully renamed %s %s to %s\n", typeName, oldName, newName)
		},
	}

	parent.AddCommand(cmd)
	return cmd
}
