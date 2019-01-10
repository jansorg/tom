package store

import (
	"fmt"
	"strings"

	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/query"
)

func NewStoreHelper(store model.Store) *Helper {
	return &Helper{
		store: store,
		query: query.NewStoreQuery(store),
	}
}

type Helper struct {
	store model.Store
	query query.StoreQuery
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

func (s *Helper) RenameProjectByName(oldName, newName string) (*model.Project, error) {
	if p, err := s.query.ProjectByFullNameOrID(oldName); err != nil {
		return nil, fmt.Errorf("no project found for '%s'", oldName)
	} else {
		return s.RenameProject(p, newName)
	}
}

func (s *Helper) RenameProject(project *model.Project, newName string) (*model.Project, error) {
	if project == nil {
		return nil, fmt.Errorf("project is undefined")
	} else if newName == "" {
		return nil, fmt.Errorf("new name is empty")
	} else if _, err := s.query.ProjectByFullNameOrID(newName); err == nil {
		return nil, fmt.Errorf("project %s already exists", newName)
	}

	// now find parent if newName indicates a nested project, just rename if it's a top-level project

	if !strings.Contains(newName, "/") {
		// now top-level
		project.ParentID = ""
		project.Name = newName
		return s.store.UpdateProject(*project)
	}

	// now a nested project
	parts := strings.Split(newName, "/")
	parentNames := parts[:len(parts)-1]
	parent, _, err := s.GetOrCreateNestedProjectNames(parentNames...)
	if err != nil {
		return nil, err
	}

	project.ParentID = parent.ID
	project.Name = parts[len(parts)-1]
	return s.store.UpdateProject(*project)
}
