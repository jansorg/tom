package query

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/slices"
	"github.com/jansorg/tom/go-tom/store"
	"github.com/jansorg/tom/go-tom/util"
)

type StoreQuery interface {
	AnyByID(id string) (interface{}, error)

	ProjectByID(id string) (*model.Project, error)
	ProjectByFullName(names []string) (*model.Project, error)
	ProjectByFullNameOrID(nameOrID string, delimiter string) (*model.Project, error)
	ProjectsByShortName(name string) []*model.Project
	ProjectsByShortNameOrID(nameOrID string) []*model.Project
	WithProjectAndParents(id string, f func(*model.Project) bool) bool
	CollectProjectAndSubprojects(id string) ([]*model.Project, error)
	CollectSubprojectIDs(id string) []string
	FindRecentlyTrackedProjects(max int) (model.ProjectList, error)
	FindSuitableProject(id string, choices []string) (string, error)

	TagByID(id string) (*model.Tag, error)
	TagByName(name string) (*model.Tag, error)
	TagsByName(names ...string) ([]*model.Tag, error)

	FrameByID(id string) (*model.Frame, error)
	FramesByID(id ...string) ([]*model.Frame, error)
	FramesByProject(id string, includeSubprojects bool) model.FrameList
	FramesByTag(id string) []*model.Frame
	ActiveFrames() []*model.Frame

	IsToplevelProject(id string) bool

	FindPropertyValue(propertyID string, projectId string) (interface{}, error)
	FindPropertyValues(projectId string) map[*model.Property]interface{}
	FindPropertyByNameOrID(nameOrID string) (*model.Property, error)
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

func (q *defaultStoreQuery) ProjectByFullNameOrID(nameOrID string, delimiter string) (*model.Project, error) {
	if project, err := q.store.ProjectByID(nameOrID); err == nil {
		return project, err
	}
	return q.ProjectByFullName(strings.Split(nameOrID, delimiter))
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

// Iterates the project and its parent hierarchy until there's no parent or the function returns false
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

// Iterates the project and its hierarchy below, it stops as soon as f returns false
func (q *defaultStoreQuery) CollectProjectAndSubprojects(id string) ([]*model.Project, error) {
	if _, err := q.store.ProjectByID(id); err != nil {
		return nil, err
	}

	var result []*model.Project
	for _, p := range q.store.Projects() {
		if q.store.ProjectIsSameOrChild(id, p.ID) {
			result = append(result, p)
		}
	}
	return result, nil
}

func (q *defaultStoreQuery) CollectSubprojectIDs(id string) []string {
	if _, err := q.store.ProjectByID(id); err != nil {
		return []string{}
	}

	var result []string
	for _, p := range q.store.Projects() {
		if p.ID != id && q.store.ProjectIsSameOrChild(id, p.ID) {
			result = append(result, p.ID)
		}
	}
	return result
}

func (q *defaultStoreQuery) FindRecentlyTrackedProjects(max int) (model.ProjectList, error) {
	// frames are always sorted, collect the max distict projects
	result := model.ProjectList{}
	projectMap := map[string]bool{}

	frames := q.store.Frames()
	for i := len(frames) - 1; i >= 0; i-- {
		frame := frames[i]
		if _, ok := projectMap[frame.ProjectId]; !ok {
			projectMap[frame.ProjectId] = true
			if project, err := q.store.ProjectByID(frame.ProjectId); err != nil {
				return nil, err
			} else {
				result = append(result, project)
				if len(result) >= max {
					break
				}
			}
		}
	}
	return result, nil
}

func (q *defaultStoreQuery) FindSuitableProject(id string, choices []string) (string, error) {
	acceptedIDs := util.MapStrings(choices)
	if acceptedIDs[id] {
		return id, nil
	}

	result := ""
	q.WithProjectAndParents(id, func(project *model.Project) bool {
		currentID := project.ID
		if acceptedIDs[currentID] {
			result = currentID
			return false
		}
		return true
	})

	if result == "" {
		return "", fmt.Errorf("no suitable project found")
	}
	return result, nil
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

func (q *defaultStoreQuery) FramesByID(ids ...string) ([]*model.Frame, error) {
	var set = map[string]bool{}
	for _, id := range ids {
		set[id] = true
	}

	var frames []*model.Frame
	for _, id := range ids {
		if f, err := q.FrameByID(id); err != nil {
			return nil, err
		} else {
			frames = append(frames, f)
		}
	}
	return frames, nil
}

func (q *defaultStoreQuery) FramesByProject(id string, includeSubprojects bool) model.FrameList {
	frames, _ := q.store.FindFrames(func(f *model.Frame) (bool, error) {
		return f.ProjectId == id || includeSubprojects && q.store.ProjectIsSameOrChild(id, f.ProjectId), nil
	})
	return frames
}

func (q *defaultStoreQuery) FramesByTag(id string) []*model.Frame {
	frames, _ := q.store.FindFrames(func(f *model.Frame) (bool, error) {
		return false, nil
	})
	return frames
}

func (q *defaultStoreQuery) ActiveFrames() []*model.Frame {
	frames, _ := q.store.FindFrames(func(f *model.Frame) (bool, error) {
		return f.IsActive(), nil
	})
	return frames
}

func (q *defaultStoreQuery) FindPropertyValue(propertyID string, projectID string) (interface{}, error) {
	var value interface{}
	found := false

	prop, err := q.store.GetProperty(propertyID)
	if err != nil {
		return nil, err
	}

	q.WithProjectAndParents(projectID, func(p *model.Project) bool {
		// break early for properties not applying to subprojects
		if p.ID != projectID && !prop.ApplyToSubprojects {
			found = false
			return false
		}

		if v, err := p.GetPropertyValue(propertyID); err == nil {
			found = true
			value = v
		}
		return !found
	})

	if !found {
		return nil, store.ErrPropertyNotFound
	}
	return value, nil
}

func (q *defaultStoreQuery) FindPropertyValues(projectID string) map[*model.Property]interface{} {
	result := make(map[*model.Property]interface{})

	for _, prop := range q.store.Properties() {
		q.WithProjectAndParents(projectID, func(p *model.Project) bool {
			// break early for properties not applying to sub projects
			if p.ID != projectID && !prop.ApplyToSubprojects {
				return false
			}

			if value, err := p.GetPropertyValue(prop.ID); err == nil {
				if _, exists := result[prop]; !exists {
					result[prop] = value
				}
			}
			return true
		})
	}

	return result
}

func (q *defaultStoreQuery) FindPropertyByNameOrID(nameOrID string) (*model.Property, error) {
	property, err := q.store.GetProperty(nameOrID)
	if err == nil {
		return property, nil
	}

	for _, prop := range q.store.Properties() {
		if prop.Name == nameOrID {
			return prop, nil
		}
	}

	return nil, store.ErrPropertyNotFound
}
