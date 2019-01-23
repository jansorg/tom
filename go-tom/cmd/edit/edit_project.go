package edit

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/cmd/util"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

func newEditProjectCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var name string
	var parentNameOrID string
	var nameDelimiter string
	// fixme add properties

	var cmd = &cobra.Command{
		Use:   "project fullName | ID",
		Short: "edit properties of a project",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("name").Changed && len(name) == 0 {
				util.Fatal("unable to use empty project name")
			} else if !cmd.Flag("name").Changed && !cmd.Flag("parent").Changed {
				util.Fatalf("no modification defined, use --name or --parent to update project data")
			}

			var parent *string
			if cmd.Flag("parent").Changed {
				parent = &(parentNameOrID)
			}

			if err := doEditProjectCommand(name, parent, nameDelimiter, args, ctx); err != nil {
				util.Fatal(err)
			} else {
				println("Successfully updated project data")
			}
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "update the project name")
	cmd.Flags().StringVarP(&parentNameOrID, "parent", "p", "", "update the parent. Use an empty ID to make it a top-level project. A project keeps all frames and subprojects when it's assigned to a new parent project.")
	cmd.Flags().StringVarP(&nameDelimiter, "name-delimiter", "", "/", "Delimiter used in full project names")

	parent.AddCommand(cmd)
	return cmd
}

func doEditProjectCommand(newName string, parentNameOrID *string, nameDelimiter string, projectIDsOrNames []string, ctx *context.GoTimeContext) error {
	var err error
	var parentProjectID string

	if parentNameOrID != nil {
		if parent, err := ctx.Query.ProjectByFullNameOrID(*parentNameOrID, nameDelimiter); err != nil {
			return fmt.Errorf("parent project %s not found", *parentNameOrID)
		} else {
			parentProjectID = parent.ID
		}
	}

	// batch mode to handle many projects at once
	ctx.Store.StartBatch()
	defer ctx.Store.StopBatch()

	var projects []*model.Project
	for _, idOrName := range projectIDsOrNames {
		var project *model.Project
		if project, err = ctx.Query.ProjectByID(idOrName); err != nil {
			if project, err = ctx.Query.ProjectByFullName(strings.Split(idOrName, nameDelimiter)); err != nil {
				util.Fatalf("project %s not found", idOrName)
			}
		}
		projects = append(projects, project)
	}

	for _, p := range projects {
		if len(newName) > 0 {
			p.Name = newName
		}

		if parentNameOrID != nil {
			p.ParentID = parentProjectID
		}

		if _, err = ctx.Store.UpdateProject(*p); err != nil {
			util.Fatalf("error updating project %s: %s", p.GetFullName(nameDelimiter), err.Error())
		}
	}

	return nil
}
