package query

import (
	"fmt"
	"sort"

	"github.com/jansorg/tom/go-tom/config"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/slices"
)

type StoreQuery interface {
	AnyByID(id string) (interface{}, error)

	ProjectByID(id string) (*model.Project, error)
	ProjectByFullName(names []string) (*model.Project, error)
	ProjectsByShortName(name string) []*model.Project
	ProjectsByShortNameOrID(nameOrID string) []*model.Project
	WithProjectAndParents(id string, f func(*model.Project) bool) bool
	GetInheritedStringProp(projectID string, prop config.StringProperty) (string, bool)
	GetInheritedFloatProp(projectID string, prop config.FloatProperty) (float64, bool)
	GetInheritedIntProp(projectID string, prop config.IntProperty) (int64, bool)

	TagByID(id string) (*model.Tag, error)
	TagByName(name string) (*model.Tag, error)
	TagsByName(names ...string) ([]*model.Tag, error)

	FrameByID(id string) (*model.Frame, error)
	FramesByProject(id string, includeSubprojects bool) []*model.Frame
	FramesByTag(id string) []*model.Frame
	ActiveFrames() []*model.Frame
	IsToplevelProject(id string) bool
}

func NewStoreQuery(store model.Store) StoreQuery {
	return &defaultStoreQuery{store: store}
}

type defaultStoreQuery struct {
	store model.Store
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

func (q *defaultStoreQuery) ProjectByID(id string) (*model.Project, error) {
	return q.store.ProjectByID(id)
}

func (q *defaultStoreQuery) IsToplevelProject(id string) bool {
	p, err := q.ProjectByID(id)
	return err != nil && p.ParentID == ""
}

func (q *defaultStoreQuery) ProjectByFullName(name []string) (*model.Project, error) {
	return q.store.FindFirstProject(func(p *model.Project) bool {
		return slices.StringsEqual(name, p.FullName)
	})
}

func (q *defaultStoreQuery) ProjectsByShortName(name string) []*model.Project {
	return q.store.FindProjects(func(project *model.Project) bool {
		return project.Name == name
	})
}

func (q *defaultStoreQuery) ProjectsByShortNameOrID(nameOrID string) []*model.Project {
	if p, err := q.ProjectByID(nameOrID); err == nil {
		return []*model.Project{p}
	}
	return q.ProjectsByShortName(nameOrID)
}

// Iterates the project and its parent hierarchy until there's not parent or the function returns false
func (q *defaultStoreQuery) WithProjectAndParents(id string, f func(project *model.Project) bool) bool {
	for id != "" {
		current, err := q.ProjectByID(id)
		if err != nil {
			return false
		}

		ok := f(current)
		if !ok {
			return false
		}

		id = current.ParentID
	}
	return false
}

func (q *defaultStoreQuery) GetInheritedStringProp(projectID string, prop config.StringProperty) (string, bool) {
	value := ""
	ok := false

	q.WithProjectAndParents(projectID, func(project *model.Project) bool {
		value, ok = prop.Get(project)
		return !ok
	})

	return value, ok
}

func (q *defaultStoreQuery) GetInheritedIntProp(projectID string, prop config.IntProperty) (int64, bool) {
	var value int64
	ok := false

	q.WithProjectAndParents(projectID, func(project *model.Project) bool {
		value, ok = prop.Get(project)
		return ok
	})

	return value, ok
}

func (q *defaultStoreQuery) GetInheritedFloatProp(projectID string, prop config.FloatProperty) (float64, bool) {
	var value float64
	ok := false

	q.WithProjectAndParents(projectID, func(project *model.Project) bool {
		value, ok = prop.Get(project)
		return ok
	})

	return value, ok
}

func (q *defaultStoreQuery) TagByID(id string) (*model.Tag, error) {
	tag, err := q.store.FindFirstTag(func(t *model.Tag) bool {
		return t.ID == id
	})
	if err != nil {
		return nil, fmt.Errorf("no tag found for id %s", id)
	}
	return tag, nil
}

func (q *defaultStoreQuery) TagByName(name string) (*model.Tag, error) {
	tag, err := q.store.FindFirstTag(func(t *model.Tag) bool {
		return t.Name == name
	})

	if err != nil {
		return nil, fmt.Errorf("no tag found for name %s", name)
	}
	return tag, nil
}

func (q *defaultStoreQuery) TagsByName(names ...string) ([]*model.Tag, error) {
	sort.Strings(names)
	matching := q.store.FindTags(func(t *model.Tag) bool {
		i := sort.SearchStrings(names, t.Name)
		return i < len(names) && names[i] == t.Name
	})

	if len(matching) != len(names) {
		return nil, fmt.Errorf("unable to find all tags for %s", names)
	}
	return matching, nil
}

func (q *defaultStoreQuery) FrameByID(id string) (*model.Frame, error) {
	return q.store.FindFirstFrame(func(f *model.Frame) bool {
		return f.ID == id
	})
}

func (q *defaultStoreQuery) FramesByProject(id string, includeSubprojects bool) []*model.Frame {
	return q.store.FindFrames(func(f *model.Frame) bool {
		return f.ProjectId == id || includeSubprojects && q.store.ProjectIsChild(id, f.ProjectId)
	})
}

func (q *defaultStoreQuery) FramesByTag(id string) []*model.Frame {
	return q.store.FindFrames(func(f *model.Frame) bool {
		return false
	})
}

func (q *defaultStoreQuery) ActiveFrames() []*model.Frame {
	return q.store.FindFrames(func(f *model.Frame) bool {
		return f.IsActive()
	})
}
