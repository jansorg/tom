package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

type projectList struct {
	projects      []*model.Project
	nameDelimiter string
}

func (o projectList) size() int {
	return len(o.projects)
}

func (o projectList) get(index int, prop string, format string) (interface{}, error) {
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
	default:
		return "", fmt.Errorf("unknown property %s", prop)
	}
}

func newProjectsCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	nameDelimiter := ""

	var cmd = &cobra.Command{
		Use:   "projects",
		Short: "Prints projects",
		Run: func(cmd *cobra.Command, args []string) {
			// fixme replace with ProjectList
			projects := ctx.Store.Projects()
			sort.SliceStable(projects, func(i, j int) bool {
				return strings.Compare(projects[i].GetFullName("/"), projects[j].GetFullName("/")) < 0
			})

			list := projectList{projects: projects, nameDelimiter: nameDelimiter}
			err := printList(cmd, list, ctx)
			if err != nil {
				Fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&nameDelimiter, "name-delimiter", "", "/", "Delimiter used in the full project name")

	addListOutputFlags(cmd, "fullName", []string{"id", "fullName", "name", "parentID", "trackedDay"})
	parent.AddCommand(cmd)
	return cmd
}
