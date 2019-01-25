package dataImport

import (
	"fmt"

	"github.com/jansorg/tom/go-tom/context"
)

type Result struct {
	CreatedProjects int
	ReusedProjects  int
	CreatedFrames   int
	CreatedTags     int
}

func (r Result) String() string {
	return fmt.Sprintf("Successfully created %d projects, %d tags, and %d frames.", r.CreatedProjects, r.CreatedTags, r.CreatedFrames)
}

type Handler interface {
	Import(filename string, ctx *context.TomContext) (Result, error)
}
