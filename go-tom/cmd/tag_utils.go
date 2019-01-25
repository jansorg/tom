package cmd

import (
	"fmt"
	"strings"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

func argsToTags(ctx *context.TomContext, args []string) ([]*model.Tag, error) {
	if len(args) == 0 {
		return []*model.Tag{}, nil
	}

	var tags []*model.Tag

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
