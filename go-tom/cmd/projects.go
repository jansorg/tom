package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/report"
)

type projectList struct {
	projects []*model.Project
	reports  map[string]*report.ProjectSummary
}

func (o projectList) size() int {
	return len(o.projects)
}

func (o projectList) get(index int, prop string, format string) (interface{}, error) {
	r := o.reports[o.projects[index].ID]
	if r == nil {
		r = &report.ProjectSummary{}
	}

	switch prop {
	case "id":
		return o.projects[index].ID, nil
	case "parentID":
		return o.projects[index].ParentID, nil
	case "fullName":
		return o.projects[index].FullName, nil
	case "name":
		return o.projects[index].Name, nil
	case "trackedDay":
		return r.TrackedDay, nil
	case "trackedWeek":
		return r.TrackedWeek, nil
	case "trackedMonth":
		return r.TrackedMonth, nil
	case "trackedYear":
		return r.TrackedYear, nil
	case "totalTrackedDay":
		return r.TotalTrackedDay, nil
	case "totalTrackedWeek":
		duration := r.TotalTrackedWeek
		return duration, nil
	case "totalTrackedMonth":
		duration := r.TotalTrackedMonth
		return duration, nil
	case "totalTrackedYear":
		return r.TotalTrackedYear, nil
	default:
		return "", fmt.Errorf("unknown property %s", prop)
	}
}

func newProjectsCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "projects",
		Short: "Prints projects",
		Run: func(cmd *cobra.Command, args []string) {
			projects := ctx.Store.Projects()
			sort.SliceStable(projects, func(i, j int) bool {
				return strings.Compare(projects[i].FullName, projects[j].FullName) < 0
			})

			// fixme create only when needed
			frames := model.NewFrameList(ctx.Store.Frames())
			projectReports := report.CreateProjectReports(frames, time.Now(), ctx)

			projectList := projectList{projects: projects, reports: projectReports}

			err := printList(cmd, projectList, ctx)
			if err != nil {
				fatal(err)
			}
		},
	}

	addListOutputFlags(cmd, "fullName", []string{"id", "fullName", "name", "parentID", "trackedDay", "trackedWeek", "trackedMonth", "trackedYear", "totalTrackedDay", "totalTrackedWeek", "totalTrackedMonth", "totalTrackedYear"})
	parent.AddCommand(cmd)
	return cmd
}
