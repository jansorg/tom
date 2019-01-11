package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/report"
)

type projectList struct {
	projects []*model.Project
	reports  map[string]report.ProjectSummary
}

func (o projectList) size() int {
	return len(o.projects)
}

func (o projectList) get(index int, prop string, format string) (string, error) {
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
		duration := o.reports[o.projects[index].ID].TrackedDay
		// if format == "json" {
		return strconv.FormatInt(duration.Nanoseconds()/1000/1000, 10), nil
		// }
		// return ctx.DurationPrinter.Short(duration), nil
	case "trackedWeek":
		duration := o.reports[o.projects[index].ID].TrackedWeek
		// if format == "json" {
		return strconv.FormatInt(duration.Nanoseconds()/1000/1000, 10), nil
		// }
		// return ctx.DurationPrinter.Short(duration), nil
	case "trackedMonth":
		duration := o.reports[o.projects[index].ID].TrackedMonth
		// if format == "json" {
		return strconv.FormatInt(duration.Nanoseconds()/1000/1000, 10), nil
		// }
		// return ctx.DurationPrinter.Short(duration), nil
	case "trackedYear":
		duration := o.reports[o.projects[index].ID].TrackedYear
		// if format = "json" {
		return strconv.FormatInt(duration.Nanoseconds()/1000/1000, 10), nil
		// }
		// return ctx.DurationPrinter.Short(duration), nil
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
			err := printList(cmd, projectList)
			if err != nil {
				fatal(err)
			}
		},
	}

	addListOutputFlags(cmd, "fullName", []string{"id", "fullName", "name", "parentID", "trackedDay", "trackedWeek", "trackedMonth", "trackedYear"})
	parent.AddCommand(cmd)
	return cmd
}
