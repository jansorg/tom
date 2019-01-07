package query

import (
	"fmt"
	"sort"

	"github.com/jansorg/tom/go-tom/config"
	"github.com/jansorg/tom/go-tom/store"
)

type StoreQuery interface {
	AnyByID(id string) (interface{}, error)

	ProjectByID(id string) (*store.Project, error)
	ProjectByFullName(name string) (*store.Project, error)
	ProjectByFullNameOrID(name string) (*store.Project, error)
	ProjectsByShortName(name string) []*store.Project
	ProjectsByShortNameOrID(nameOrID string) []*store.Project
	WithProjectAndParents(id string, f func(*store.Project) bool) bool
	GetInheritedStringProp(projectID string, prop config.StringProperty) (string, bool)
	GetInheritedFloatProp(projectID string, prop config.FloatProperty) (float64, bool)
	GetInheritedIntProp(projectID string, prop config.IntProperty) (int64, bool)

	TagByID(id string) (*store.Tag, error)
	TagByName(name string) (*store.Tag, error)
	TagsByName(names ...string) ([]*store.Tag, error)

	FrameByID(id string) (*store.Frame, error)
	FramesByProject(id string, includeSubprojects bool) []*store.Frame
	FramesByTag(id string) []*store.Frame
	ActiveFrames() []*store.Frame
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

func (q *defaultStoreQuery) ProjectByFullNameOrID(nameOrID string) (*store.Project, error) {
	if p, err := q.ProjectByID(nameOrID); err == nil {
		return p, nil
	}

	return q.store.FindFirstProject(func(p *store.Project) bool {
		return p.FullName == nameOrID
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

// Iterates the project and its parent hierarchy until there's not parent or the function returns false
func (q *defaultStoreQuery) WithProjectAndParents(id string, f func(project *store.Project) bool) bool {
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

	q.WithProjectAndParents(projectID, func(project *store.Project) bool {
		value, ok = prop.Get(project)
		return !ok
	})

	return value, ok
}

func (q *defaultStoreQuery) GetInheritedIntProp(projectID string, prop config.IntProperty) (int64, bool) {
	var value int64
	ok := false

	q.WithProjectAndParents(projectID, func(project *store.Project) bool {
		value, ok = prop.Get(project)
		return ok
	})

	return value, ok
}

func (q *defaultStoreQuery) GetInheritedFloatProp(projectID string, prop config.FloatProperty) (float64, bool) {
	var value float64
	ok := false

	q.WithProjectAndParents(projectID, func(project *store.Project) bool {
		value, ok = prop.Get(project)
		return ok
	})

	return value, ok
}

func (q *defaultStoreQuery) TagByID(id string) (*store.Tag, error) {
	tag, err := q.store.FindFirstTag(func(t *store.Tag) bool {
		return t.ID == id
	})
	if err != nil {
		return nil, fmt.Errorf("no tag found for id %s", id)
	}
	return tag, nil
}

func (q *defaultStoreQuery) TagByName(name string) (*store.Tag, error) {
	tag, err := q.store.FindFirstTag(func(t *store.Tag) bool {
		return t.Name == name
	})

	if err != nil {
		return nil, fmt.Errorf("no tag found for name %s", name)
	}
	return tag, nil
}

func (q *defaultStoreQuery) TagsByName(names ...string) ([]*store.Tag, error) {
	sort.Strings(names)
	matching := q.store.FindTags(func(t *store.Tag) bool {
		i := sort.SearchStrings(names, t.Name)
		return i < len(names) && names[i] == t.Name
	})

	if len(matching) != len(names) {
		return nil, fmt.Errorf("unable to find all tags for %s", names)
	}
	return matching, nil
}

func (q *defaultStoreQuery) FrameByID(id string) (*store.Frame, error) {
	return q.store.FindFirstFrame(func(f *store.Frame) bool {
		return f.ID == id
	})
}

func (q *defaultStoreQuery) FramesByProject(id string, includeSubprojects bool) []*store.Frame {
	return q.store.FindFrames(func(f *store.Frame) bool {
		return f.ProjectId == id || includeSubprojects && q.store.ProjectIsChild(id, f.ProjectId)
	})
}

func (q *defaultStoreQuery) FramesByTag(id string) []*store.Frame {
	return q.store.FindFrames(func(f *store.Frame) bool {
		return false
	})
}

func (q *defaultStoreQuery) ActiveFrames() []*store.Frame {
	return q.store.FindFrames(func(f *store.Frame) bool {
		return f.IsActive()
	})
}
