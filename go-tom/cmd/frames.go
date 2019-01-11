package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

type frameList []*model.Frame

func (f frameList) size() int {
	return len(f)
}

func (f frameList) get(index int, prop string, format string) (interface{}, error) {
	switch prop {
	case "id":
		return f[index].ID, nil
	case "projectName":
		if project, err := ctx.Query.ProjectByID(f[index].ProjectId); err != nil {
			return "", err
		} else {
			return project.Name, nil
		}
	case "projectFullName":
		if project, err := ctx.Query.ProjectByID(f[index].ProjectId); err != nil {
			return "", err
		} else {
			return project.FullName, nil
		}
	case "startTime":
		return f[index].Start.In(time.Local), nil
	case "stopTime":
		frame := f[index]
		if frame.IsActive() {
			return "", nil
		}
		return frame.End.In(time.Local), nil
	case "lastUpdated":
		frame := f[index]
		if frame.Updated == nil {
			return "", nil
		}
		return frame.Updated, nil;
	case "duration":
		frame := f[index]
		return frame.Duration(), nil
	case "notes":
		frame := f[index]
		return frame.Notes, nil;
	case "tagIDs":
		frame := f[index]
		return strings.Join(frame.TagIDs, ","), nil;
	default:
		return "", fmt.Errorf("unknown property %s", prop)
	}
}

func newFramesCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	projectIDOrName := ""
	includeSubprojects := false

	var cmd = &cobra.Command{
		Use:   "frames",
		Short: "Print a listing of all frames",
		Run: func(cmd *cobra.Command, args []string) {
			var frames frameList
			if projectIDOrName == "" {
				frames = ctx.Store.Frames()
			} else {
				project, err := ctx.Query.ProjectByFullNameOrID(projectIDOrName)
				if err != nil {
					fatal(fmt.Errorf("no project found for %s", projectIDOrName))
				}
				frames = ctx.Query.FramesByProject(project.ID, includeSubprojects)
			}

			if err := printList(cmd, frames, ctx); err != nil {
				fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&projectIDOrName, "project", "p", "", "Only frames of this project will be printed. Project IDs or full project names are accepted. Default: no project")
	cmd.Flags().BoolVarP(&includeSubprojects, "subprojects", "s", false, "Include frames of subprojects")
	addListOutputFlags(cmd, "id", []string{"id", "projectID", "projectName", "projectFullName", "startTime", "stopTime", "duration", "lastUpdated", "notes", "tagIDs"})

	parent.AddCommand(cmd)
	return cmd
}
