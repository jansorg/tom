package query

import (
	"fmt"

	"github.com/jansorg/gotime/gotime/store"
)

type StoreQuery interface {
	AnyByID(id string) (interface{}, error)

	ProjectByID(id string) (*store.Project, error)
	ProjectByFullName(name string) (*store.Project, error)
	ProjectsByShortName(name string) []*store.Project
	ProjectsByShortNameOrID(nameOrID string) []*store.Project

	TagByID(id string) (*store.Tag, error)
	TagByName(name string) (*store.Tag, error)

	FrameByID(id string) (*store.Frame, error)
	FramesByProject(id string) []*store.Frame
	FramesByTag(id string) []*store.Frame
	UnstoppedFrames() []*store.Frame
	IsToplevelProject(id string) bool
}

func NewStoreQuery(store store.Store) StoreQuery {
	return &defaultStoreQuery{store: store}
}

type defaultStoreQuery struct {
	store store.Store
}

func (q *defaultStoreQuery) AnyByID(id string) (interface{}, error) {
	var v interface{}
	var err error

	if v, err = q.ProjectByID(id); err == nil {
		return v, nil
	}

	if v, err = q.TagByID(id); err == nil {
		return v, nil
	}

	if v, err = q.FrameByID(id); err == nil {
		return v, nil
	}

	return nil, fmt.Errorf("no data found for id %s", id)
}

func (q *defaultStoreQuery) ProjectByID(id string) (*store.Project, error) {
	return q.store.ProjectByID(id)
}

func (q *defaultStoreQuery) IsToplevelProject(id string) bool {
	p, err := q.ProjectByID(id)
	return err != nil && p.ParentID == ""
}

func (q *defaultStoreQuery) ProjectByFullName(name string) (*store.Project, error) {
	return q.store.FindFirstProject(func(p *store.Project) bool {
		return p.FullName == name
	})
}

func (q *defaultStoreQuery) ProjectsByShortName(name string) []*store.Project {
	return q.store.FindProjects(func(project *store.Project) bool {
		return project.Name == name
	})
}

func (q *defaultStoreQuery) ProjectsByShortNameOrID(nameOrID string) []*store.Project {
	return q.store.FindProjects(func(p *store.Project) bool {
		return p.ID == nameOrID || p.FullName == nameOrID
	})
}

func (q *defaultStoreQuery) TagByID(id string) (*store.Tag, error) {
	return q.store.FindFirstTag(func(t *store.Tag) bool {
		return t.ID == id
	})
}

func (q *defaultStoreQuery) TagByName(name string) (*store.Tag, error) {
	return q.store.FindFirstTag(func(t *store.Tag) bool {
		return t.Name == name
	})
}

func (q *defaultStoreQuery) FrameByID(id string) (*store.Frame, error) {
	return q.store.FindFirstFrame(func(f *store.Frame) bool {
		return f.ID == id
	})
}

func (q *defaultStoreQuery) FramesByProject(id string) []*store.Frame {
	return q.store.FindFrames(func(f *store.Frame) bool {
		return f.ProjectId == id
	})
}

func (q *defaultStoreQuery) FramesByTag(id string) []*store.Frame {
	return q.store.FindFrames(func(f *store.Frame) bool {
		return false
	})
}

func (q *defaultStoreQuery) UnstoppedFrames() []*store.Frame {
	return q.store.FindFrames(func(f *store.Frame) bool {
		return f.IsActive()
	})
}
