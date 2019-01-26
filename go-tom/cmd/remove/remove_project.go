package remove

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/cmd/util"
	"github.com/jansorg/tom/go-tom/context"
)

func newRemoveProjectCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	nameDelimiter := ""

	var cmd = &cobra.Command{
		Use:   "project <project name or project ID> ...",
		Short: "removes new project and all its associated data, including subprojects and time entries",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if removedProjects, removedFrames, err := doRemoveProjects(ctx, nameDelimiter, args); err != nil {
				util.Fatal("Error removing projects: %s", err.Error())
			} else {
				fmt.Printf("Successfully removed %d projects and %d frames\n", removedProjects, removedFrames)
			}
		},
	}

	cmd.Flags().StringVarP(&nameDelimiter, "name-delimiter", "", "/", "Delimiter used in full project names")
	parent.AddCommand(cmd)
	return cmd
}

func doRemoveProjects(ctx *context.TomContext, nameDelimiter string, idOrNames []string) (removedProjects, removedFrames int, err error) {
	ctx.Store.StartBatch()
	defer ctx.Store.StopBatch()

	for _, idOrName := range idOrNames {
		projects, err := ctx.Query.ProjectByFullNameOrID(idOrName, nameDelimiter)
		if err != nil {
			return 0, 0, err
		}

		projectCount, frameCount, err := ctx.StoreHelper.RemoveProject(projects)
		if err != nil {
			return 0, 0, err
		}

		removedProjects += projectCount
		removedFrames += frameCount
	}
	return
}
