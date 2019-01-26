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

func (s *Helper) RenameProjectByIDOrName(oldName, newName string) (*model.Project, error) {
	if strings.TrimSpace(newName) == "" {
		return nil, fmt.Errorf("new project name must not be empty")
	}

	p, err := s.query.ProjectByID(oldName)
	if err != nil {
		p, err = s.query.ProjectByFullName(strings.Split(oldName, "/"))
		if err != nil {
			return nil, fmt.Errorf("no project found for '%s'", oldName)
		}
	}

	return s.RenameProject(p, strings.Split(newName, "/"), true)
}

func (s *Helper) RenameProject(project *model.Project, newName []string, allowHierarchyUpdate bool) (*model.Project, error) {
	if project == nil {
		return nil, fmt.Errorf("project is undefined")
	} else if len(newName) == 0 {
		return nil, fmt.Errorf("new name is empty")
	} else if _, err := s.query.ProjectByFullName(newName); err == nil {
		return nil, fmt.Errorf("project %s already exists", newName)
	}

	if !allowHierarchyUpdate {
		if len(newName) != 1 {
			return nil, fmt.Errorf("project rename without hierarchy update needs just one name element")
		}
		project.Name = newName[0]
		return s.store.UpdateProject(*project)
	}

	// now find parent if newName indicates a nested project, just rename if it's a top-level project
	if len(newName) == 1 {
		// now top-level
		project.ParentID = ""
		project.Name = newName[0]
		return s.store.UpdateProject(*project)
	}

	// now a nested project
	parentNames := newName[:len(newName)-1]
	parent, _, err := s.GetOrCreateNestedProjectNames(parentNames...)
	if err != nil {
		return nil, err
	}

	project.ParentID = parent.ID
	project.Name = newName[len(newName)-1]
	return s.store.UpdateProject(*project)
}

func (s *Helper) MoveProject(project *model.Project, newParentID string) (*model.Project, error) {
	if newParentID != "" {
		// return an error if moved onto itself or into own child scope
		reject := false
		s.query.WithProjectAndParents(newParentID, func(parent *model.Project) bool {
			reject = reject || parent.ID == project.ID
			return !reject
		})

		if reject {
			return nil, fmt.Errorf("moving a project into its own child scope is not allowed")
		}
	}

	project.ParentID = newParentID
	return project, nil
}

func (s *Helper) RemoveProject(project *model.Project) (int, int, error) {
	if project == nil {
		return 0, 0, fmt.Errorf("project undefined")
	}

	// collect all sub projects
	projects, err := s.query.CollectProjectAndSubprojects(project.ID)
	if err != nil {
		return 0, 0, err
	}

	removedProjects := 0
	removedFrames := 0
	for _, p := range projects {
		frames := s.query.FramesByProject(p.ID, true)
		for _, f := range frames {
			if err := s.store.RemoveFrame(f.ID); err != nil {
				return 0, 0, err
			}
			removedFrames++
		}

		if err := s.store.RemoveProject(p.ID); err != nil {
			return 0, 0, err
		}
		removedProjects++
	}

	return removedProjects, removedFrames, nil
}
