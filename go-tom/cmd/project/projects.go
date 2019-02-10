package project

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/cmd/cmdUtil"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/util"
)

type projectList struct {
	projects      []*model.Project
	nameDelimiter string
}

func (o projectList) Size() int {
	return len(o.projects)
}

func (o projectList) Get(index int, prop string, format string, ctx *context.TomContext) (interface{}, error) {
	switch prop {
	case "id":
		return o.projects[index].ID, nil
	case "parentID":
		return o.projects[index].ParentID, nil
	case "fullName":
		if format == "json" {
			return o.projects[index].FullName, nil
		}
		return o.projects[index].GetFullName(o.nameDelimiter), nil
	case "name":
		return o.projects[index].Name, nil
	case "hourlyRate":
		rate := o.projects[index].HourlyRate()
		if rate == nil {
			return "", nil
		}
		return rate.ParsableString(), nil
	case "appliedHourlyRate":
		rate, err := ctx.Query.HourlyRate(o.projects[index].ID)
		if rate == nil || err != nil {
			return "", nil
		}
		return rate.ParsableString(), nil
	default:
		return "", fmt.Errorf("unknown property %s", prop)
	}
}

func NewCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	nameDelimiter := ""
	recentProjects := 0

	var cmd = &cobra.Command{
		Use:   "projects",
		Short: "Prints projects",
		Run: func(cmd *cobra.Command, args []string) {
			var projects model.ProjectList
			if cmd.Flag("recent").Changed {
				var err error
				if projects, err = ctx.Query.FindRecentlyTrackedProjects(recentProjects); err != nil {
					log.Fatal(err)
				}
			} else {
				projects = ctx.Store.Projects()
				projects.SortByFullname()
			}

			list := projectList{projects: projects, nameDelimiter: nameDelimiter}
			err := cmdUtil.PrintList(cmd, list, ctx)
			if err != nil {
				util.Fatal(err)
			}
		},
	}

	cmd.Flags().IntVarP(&recentProjects, "recent", "", 0, "If set then only the most recently tracked projects will be returned.")
	cmd.Flags().StringVarP(&nameDelimiter, "name-delimiter", "", "/", "Delimiter used in the full project name")
	cmdUtil.AddListOutputFlags(cmd, "fullName", []string{"id", "fullName", "name", "parentID", "hourlyRate", "appliedHourlyRate"})

	parent.AddCommand(cmd)
	return cmd
}
