package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/store"
)

type tagList []*store.Tag

func (o tagList) size() int {
	return len(o)
}

func (t tagList) get(index int, prop string) (string, error) {
	switch prop {
	case "id":
		return t[index].ID, nil
	case "name":
		return t[index].Name, nil
	default:
		return "", fmt.Errorf("unknown property %s", prop)
	}
}

func newTagsCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "tags",
		Short: "Prints tags",
		Run: func(cmd *cobra.Command, args []string) {
			var tags tagList = ctx.Store.Tags()
			sort.SliceStable(tags, func(i, j int) bool {
				return strings.Compare(tags[i].Name, tags[j].Name) < 0
			})

			if err := printList(cmd, tags); err != nil {
				fatal(err)
			}
		},
	}

	addListOutputFlags(cmd, "name", []string{"id", "name"})
	parent.AddCommand(cmd)
	return cmd
}
