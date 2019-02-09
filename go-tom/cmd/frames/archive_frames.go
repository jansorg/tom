package frames

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/util"
)

func newArchiveCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	projectIDOrName := ""
	nameDelimiter := ""
	includeSubprojects := false

	var cmd = &cobra.Command{
		Use:   "archive",
		Short: "Archive a set of frames",
		Run: func(cmd *cobra.Command, args []string) {
			if err := archiveFrames(projectIDOrName, nameDelimiter, includeSubprojects, ctx); err != nil {
				util.Fatalf("Error archiving frames: %s", err.Error())
			} else {
				fmt.Println("archived project frames")
			}
		},
	}

	cmd.Flags().StringVarP(&projectIDOrName, "project", "p", "", "Only frames of this project will be archived")
	cmd.Flags().StringVarP(&nameDelimiter, "name-delimiter", "", "/", "Delimiter used in the full project name")
	cmd.Flags().BoolVarP(&includeSubprojects, "include-subprojects", "", includeSubprojects, "Also archive frames of all subprojects")

	parent.AddCommand(cmd)
	return cmd
}

func archiveFrames(projectIDOrName string, nameDelimiter string, includeSubprojects bool, ctx *context.TomContext) error {
	ctx.Store.StartBatch()
	defer ctx.Store.StopBatch()

	project, err := ctx.Query.ProjectByFullNameOrID(projectIDOrName, nameDelimiter)
	if err != nil {
		return err
	}

	frames := ctx.Query.FramesByProject(project.ID, includeSubprojects)
	frames.ExcludeArchived()

	for _, f := range frames {
		f.Archived = true
		_, err := ctx.Store.UpdateFrame(*f)
		if err != nil {
			return err
		}
	}

	return nil
}
