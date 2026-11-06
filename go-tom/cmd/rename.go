package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/util"
)

func newRenameCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rename TYPE {name | ID} newName",
		Short: "rename a project or a tag",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var typeName = args[0]
			var idOrOldName = args[1]
			var newName = args[2]

			switch typeName {
			case "project":
				var project *model.Project
				var err error
				if project, err = ctx.Query.ProjectByID(idOrOldName); err != nil {
					project, err = ctx.Query.ProjectByFullName(strings.Split(idOrOldName, "/"))
				}
				if err != nil {
					util.Fatalf("project %s not found", idOrOldName)
				}

				if strings.Contains(newName, "/") {
					util.Fatalf("the new name %s must not contain '/'. Use 'tom edit project --parent' to move a project", newName)
				}

				if _, err := ctx.StoreHelper.RenameProject(project, []string{newName}, false); err != nil {
					util.Fatalf("rename failed: %s", err.Error())
				}
			case "tag":
				if tag, err := ctx.Query.TagByName(idOrOldName); err != nil {
					util.Fatalf("tag %s not found", idOrOldName)
				} else {
					tag.Name = newName
					if _, err := ctx.Store.UpdateTag(*tag); err != nil {
						util.Fatalf("unable to rename tag %s to %s: %v", idOrOldName, newName, err)
					}
				}
			default:
				util.Fatal(fmt.Errorf("unknown TYPE %s. Valid values are project, tag", typeName))
			}

			fmt.Printf("successfully renamed %s %s to %s\n", typeName, idOrOldName, newName)
		},
	}

	parent.AddCommand(cmd)
	return cmd
}
