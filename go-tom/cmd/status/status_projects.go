package status

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/cmd/cmdUtil"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/report"
	"github.com/jansorg/tom/go-tom/util"
)

type projectStatusList struct {
	nameDelimiter string
	reports       []*report.ProjectSummary
}

func (o projectStatusList) Size() int {
	return len(o.reports)
}

func (o projectStatusList) Get(index int, prop string, format string, ctx *context.TomContext) (interface{}, error) {
	summary := o.reports[index]

	switch prop {
	case "id":
		return summary.Project.ID, nil
	case "parentID":
		return summary.Project.ParentID, nil
	case "fullName":
		return summary.Project.GetFullName(o.nameDelimiter), nil
	case "name":
		return summary.Project.Name, nil
	case "trackedDay":
		return summary.TrackedDay.Get(), nil
	case "trackedWeek":
		return summary.TrackedWeek.Get(), nil
	case "trackedMonth":
		return summary.TrackedMonth.Get(), nil
	case "trackedYear":
		return summary.TrackedYear.Get(), nil
	case "trackedAll":
		return summary.TrackedAll.Get(), nil
	case "totalTrackedDay":
		return summary.TrackedTotalDay.Get(), nil
	case "totalTrackedWeek":
		return summary.TrackedTotalWeek.Get(), nil
	case "totalTrackedMonth":
		return summary.TrackedTotalMonth.Get(), nil
	case "totalTrackedYear":
		return summary.TrackedTotalYear.Get(), nil
	case "totalTrackedAll":
		return summary.TrackedTotalAll.Get(), nil
	default:
		return "", fmt.Errorf("unknown property %s", prop)
	}
}

func newProjectsStatusCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	showEmpty := false
	includeActiveFrames := false
	includeArchivedFrames := true
	showOverall := false
	nameDelimiter := ""

	var cmd = &cobra.Command{
		Use:   "projects",
		Short: "Prints project status",
		Run: func(cmd *cobra.Command, args []string) {
			refTime := time.Now()

			var refEnd *time.Time
			if includeActiveFrames {
				refEnd = &refTime
			}

			projectReports := report.CreateProjectReports(refTime, showEmpty, includeArchivedFrames, refEnd, "ALL", ctx)

			var reportList []*report.ProjectSummary
			for _, v := range projectReports {
				reportList = append(reportList, v)
			}
			sort.Slice(reportList, func(i, j int) bool {
				return strings.Compare(reportList[i].Project.GetFullName("/"), reportList[j].Project.GetFullName("/")) < 0
			})

			if err := cmdUtil.PrintList(cmd, projectStatusList{reports: reportList, nameDelimiter: nameDelimiter}, ctx); err != nil {
				util.Fatal(err)
			}
		},
	}

	cmd.Flags().BoolVarP(&showEmpty, "show-empty", "e", showEmpty, "Include projects without tracked time in the output")
	cmd.Flags().BoolVarP(&showOverall, "show-overall", "", showOverall, "Show a summary of all projects, e.g. overall today. The used project ID is 'ALL'.")
	cmd.Flags().StringVarP(&nameDelimiter, "name-delimiter", "", "/", "Delimiter used in the full project name")
	cmd.Flags().BoolVarP(&includeActiveFrames, "include-active", "", includeActiveFrames, "Include active frames in the status. The current time will be used as end time of these frames.")
	cmd.Flags().BoolVarP(&includeArchivedFrames, "archived", "", includeArchivedFrames, "Include archived frames in the status.")

	cmdUtil.AddListOutputFlags(cmd, "fullName,trackedDay,trackedWeek,trackedMonth", []string{
		"id", "fullName", "name", "parentID",
		"trackedDay", "trackedWeek", "trackedMonth", "trackedYear", "trackedAll",
		"totalTrackedDay", "totalTrackedWeek", "totalTrackedMonth", "totalTrackedYear", "totalTrackedAll"})
	parent.AddCommand(cmd)
	return cmd
}
