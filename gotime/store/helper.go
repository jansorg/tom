package store

import (
	"fmt"
	"strings"
)

func NewStoreHelper(store Store) *StoreHelper {
	return &StoreHelper{store: store}
}

type StoreHelper struct {
	store Store
}

func (s *StoreHelper) GetOrCreateNestedProject(fullName string) (*Project, error) {
	return s.GetOrCreateNestedProjectNames(strings.Split(fullName, "/")...)
}

func (s *StoreHelper) GetOrCreateNestedProjectNames(names ...string) (*Project, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("empty project name")
	}

	parentID := ""
	var project *Project
	var err error

	for _, name := range names {
		project, err = s.GetOrCreateProject(name, parentID)
		if err != nil {
			return nil, err
		}

		project.ParentID = parentID
		parentID = project.ID
	}

	if project == nil {
		return nil, fmt.Errorf("unable to create project %v", names)
	}
	return project, nil
}

func (s *StoreHelper) GetOrCreateProject(shortName string, parentID string) (*Project, error) {
	// fmt.Printf("locating project %s\n", shortName)
	project, err := s.store.FindFirstProject(func(project *Project) bool {
		return project.ParentID == parentID && project.Name == shortName
	})

	if err == nil {
		return project, nil
	}

	return s.store.AddProject(Project{
		Name:     shortName,
		ParentID: parentID,
	})
}
