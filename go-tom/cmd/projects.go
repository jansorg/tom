package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

type projectList []*model.Project

func (o projectList) size() int {
	return len(o)
}

func (o projectList) get(index int, prop string, format string) (interface{}, error) {
	switch prop {
	case "id":
		return o[index].ID, nil
	case "parentID":
		return o[index].ParentID, nil
	case "fullName":
		return o[index].FullName, nil
	case "name":
		return o[index].Name, nil
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

			var projectList projectList = projects
			err := printList(cmd, projectList, ctx)
			if err != nil {
				fatal(err)
			}
		},
	}

	addListOutputFlags(cmd, "fullName", []string{"id", "fullName", "name", "parentID", "trackedDay"})
	parent.AddCommand(cmd)
	return cmd
}
