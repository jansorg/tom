package cmd

import (
	"fmt"
	"strings"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/store"
)

func argsToTags(ctx *context.GoTimeContext, args []string) ([]*store.Tag, error) {
	if len(args) == 0 {
		return []*store.Tag{}, nil
	}

	var tags []*store.Tag

	for _, arg := range args {
		if !strings.HasPrefix(arg, "+") {
			return nil, fmt.Errorf("%s is not matching +tagname", arg)
		}

		tag, err := ctx.Query.TagByName(arg[1:])
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
