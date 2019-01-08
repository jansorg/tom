package store

import (
	"fmt"
	"strings"

	"github.com/jansorg/tom/go-tom/model"
)

func NewStoreHelper(store model.Store) *Helper {
	return &Helper{store: store}
}

type Helper struct {
	store model.Store
}

func (s *Helper) GetOrCreateTag(name string) (*model.Tag, bool, error) {
	existing, err := s.store.FindFirstTag(func(tag *model.Tag) bool {
		return tag.Name == name
	});
	if err == nil {
		return existing, false, nil
	}

	tag, err := s.store.AddTag(model.Tag{Name: name})
	return tag, true, err
}

func (s *Helper) GetOrCreateNestedProject(fullName string) (*model.Project, bool, error) {
	return s.GetOrCreateNestedProjectNames(strings.Split(fullName, "/")...)
}

func (s *Helper) GetOrCreateNestedProjectNames(names ...string) (*model.Project, bool, error) {
	if len(names) == 0 {
		return nil, false, fmt.Errorf("empty project name")
	}

	parentID := ""
	var project *model.Project
	var err error

	created := false
	for _, name := range names {
		project, created, err = s.GetOrCreateProject(name, parentID)
		if err != nil {
			return nil, created, err
		}

		project.ParentID = parentID
		parentID = project.ID
	}

	if project == nil {
		return nil, false, fmt.Errorf("unable to create project %v", names)
	}
	return project, created, nil
}

func (s *Helper) GetOrCreateProject(shortName string, parentID string) (*model.Project, bool, error) {
	// fmt.Printf("locating project %s\n", shortName)
	project, err := s.store.FindFirstProject(func(project *model.Project) bool {
		return project.ParentID == parentID && project.Name == shortName
	})

	if err == nil {
		return project, false, nil
	}

	project, err = s.store.AddProject(model.Project{
		Name:     shortName,
		ParentID: parentID,
	})
	return project, true, err
}
