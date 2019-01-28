package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/cmd/util"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

type tagList []*model.Tag

func (o tagList) Size() int {
	return len(o)
}

func (t tagList) Get(index int, prop string, format string) (interface{}, error) {
	switch prop {
	case "id":
		return t[index].ID, nil
	case "name":
		return t[index].Name, nil
	default:
		return "", fmt.Errorf("unknown property %s", prop)
	}
}

func newTagsCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "tags",
		Short: "Prints tags",
		Run: func(cmd *cobra.Command, args []string) {
			var tags tagList = ctx.Store.Tags()
			sort.SliceStable(tags, func(i, j int) bool {
				return strings.Compare(tags[i].Name, tags[j].Name) < 0
			})

			if err := util.PrintList(cmd, tags, ctx); err != nil {
				util.Fatal(err)
			}
		},
	}

	util.AddListOutputFlags(cmd, "name", []string{"id", "name"})
	parent.AddCommand(cmd)
	return cmd
}
