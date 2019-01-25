package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/cmd/util"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/report"
)

type projectStatusList struct {
	nameDelimiter string
	reports       []*report.ProjectSummary
}

func (o projectStatusList) size() int {
	return len(o.reports)
}

func (o projectStatusList) get(index int, prop string, format string) (interface{}, error) {
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
		return summary.TrackedDay, nil
	case "trackedWeek":
		return summary.TrackedWeek, nil
	case "trackedMonth":
		return summary.TrackedMonth, nil
	case "trackedYear":
		return summary.TrackedYear, nil
	case "trackedAll":
		return summary.TrackedAll, nil
	case "totalTrackedDay":
		return summary.TotalTrackedDay, nil
	case "totalTrackedWeek":
		return summary.TotalTrackedWeek, nil
	case "totalTrackedMonth":
		return summary.TotalTrackedMonth, nil
	case "totalTrackedYear":
		return summary.TotalTrackedYear, nil
	case "totalTrackedAll":
		return summary.TotalTrackedAll, nil
	default:
		return "", fmt.Errorf("unknown property %s", prop)
	}
}

func newProjectsStatusCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	showEmpty := false
	showOverall := false
	nameDelimiter := ""

	var cmd = &cobra.Command{
		Use:   "projects",
		Short: "Prints project status",
		Run: func(cmd *cobra.Command, args []string) {
			projectReports := report.CreateProjectReports(time.Now(), showEmpty, "ALL", ctx)

			var reportList []*report.ProjectSummary
			for _, v := range projectReports {
				reportList = append(reportList, v)
			}
			sort.Slice(reportList, func(i, j int) bool {
				return strings.Compare(reportList[i].Project.GetFullName("/"), reportList[j].Project.GetFullName("/")) < 0
			})

			if err := printList(cmd, projectStatusList{reports: reportList, nameDelimiter: nameDelimiter}, ctx); err != nil {
				util.Fatal(err)
			}
		},
	}

	cmd.Flags().BoolVarP(&showEmpty, "show-empty", "e", showEmpty, "Includes projects without tracked time in the output")
	cmd.Flags().BoolVarP(&showOverall, "show-overall", "", showOverall, "Show a summary of all projects, e.g. overall today. The used project ID is 'ALL'.")
	cmd.Flags().StringVarP(&nameDelimiter, "name-delimiter", "", "/", "Delimiter used in the full project name")

	addListOutputFlags(cmd, "fullName,trackedDay,trackedWeek,trackedMonth", []string{
		"id", "fullName", "name", "parentID",
		"trackedDay", "trackedWeek", "trackedMonth", "trackedYear", "trackedAll",
		"totalTrackedDay", "totalTrackedWeek", "totalTrackedMonth", "totalTrackedYear", "totalTrackedAll"})
	parent.AddCommand(cmd)
	return cmd
}
