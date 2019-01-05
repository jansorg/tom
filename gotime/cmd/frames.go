package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/store"
)

type frameList []*store.Frame

func (f frameList) size() int {
	return len(f)
}

func (f frameList) get(index int, prop string) (string, error) {
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
		return f[index].Start.In(time.Local).String(), nil
	case "stopTime":
		frame := f[index]
		if frame.IsActive() {
			return "", nil
		}
		return frame.End.In(time.Local).String(), nil
	case "duration":
		frame := f[index]
		return ctx.DurationPrinter.Short(frame.Duration()), nil
	default:
		return "", fmt.Errorf("unknown property %s", prop)
	}
}

func newFramesCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	projectIDOrName := ""

	var cmd = &cobra.Command{
		Use:   "frames",
		Short: "Print a listing of all frames",
		Run: func(cmd *cobra.Command, args []string) {
			var frames frameList
			if projectIDOrName == "" {
				frames = context.Store.Frames()
			} else {
				project, err := context.Query.ProjectByFullNameOrID(projectIDOrName)
				if err != nil {
					fatal(fmt.Errorf("no project found for %s", projectIDOrName))
				}
				frames = context.Query.FramesByProject(project.ID)
			}

			if err := printList(cmd, frames); err != nil {
				fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&projectIDOrName, "project", "p", "", "Only frames of this project will be printed. Project IDs or full project names are accepted. Default: no project")
	addListOutputFlags(cmd, "name", []string{"id", "projectID", "projectName", "projectFullTime", "startTime", "stopTime", "duration"})

	parent.AddCommand(cmd)
	return cmd
}
