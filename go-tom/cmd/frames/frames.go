package frames

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/cmd/cmdUtil"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/util"
)

type frameList model.FrameList

func (f frameList) Size() int {
	return len(f)
}

func (f frameList) Get(index int, prop string, format string, ctx *context.TomContext) (interface{}, error) {
	switch prop {
	case "id":
		return f[index].ID, nil
	case "projectName":
		if project, err := ctx.Query.ProjectByID(f[index].ProjectId); err != nil {
			return "", err
		} else {
			return project.Name, nil
		}
	case "projectID":
		return f[index].ProjectId, nil
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
	case "archived":
		frame := f[index]
		return frame.Archived, nil
	default:
		return "", fmt.Errorf("unknown property %s", prop)
	}
}

func NewCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	projectIDOrName := ""
	includeSubprojects := false
	showArchived := true

	var cmd = &cobra.Command{
		Use:   "frames",
		Short: "Print a listing of all frames",
		Run: func(cmd *cobra.Command, args []string) {
			frames := filterFrames(projectIDOrName, ctx, includeSubprojects, showArchived)
			if err := cmdUtil.PrintList(cmd, frames, ctx); err != nil {
				util.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&projectIDOrName, "project", "p", "", "Only frames of this project will be printed. Project IDs or full project names are accepted. Default: no project")
	cmd.Flags().BoolVarP(&includeSubprojects, "subprojects", "s", false, "Include frames of subprojects")
	cmd.Flags().BoolVarP(&showArchived, "archived", "", showArchived, "Show/Hide archived frames")
	cmdUtil.AddListOutputFlags(cmd, "id", []string{"id", "projectID", "projectName", "projectFullName", "startTime", "stopTime", "duration", "lastUpdated", "notes", "tagIDs", "archived"})

	parent.AddCommand(cmd)
	return cmd
}

func filterFrames(projectIDOrName string, ctx *context.TomContext, includeSubprojects bool, showArchived bool) frameList {
	var frames model.FrameList
	if projectIDOrName == "" {
		frames = ctx.Store.Frames()
	} else {
		project, err := ctx.Query.ProjectByID(projectIDOrName)
		if err != nil {
			project, err = ctx.Query.ProjectByFullName(strings.Split(projectIDOrName, "/"))
			if err != nil {
				util.Fatal(fmt.Errorf("no project found for %s", projectIDOrName))
			}
		}
		frames = ctx.Query.FramesByProject(project.ID, includeSubprojects)
	}

	if !showArchived {
		frames.ExcludeArchived()
	}

	return frameList(frames)
}
