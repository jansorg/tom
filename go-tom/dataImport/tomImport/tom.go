package tomImport

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dataImport"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/store"
)

func NewImporter() dataImport.Handler {
	return &tomImporter{}
}

type tomImporter struct{}

func (tomImporter) Import(directoryPath string, ctx *context.TomContext) (dataImport.Result, error) {
	if status, err := os.Stat(directoryPath); err != nil || !status.IsDir() {
		return dataImport.Result{}, fmt.Errorf("import path does not exist or is not a directory: %s", directoryPath)
	}

	importStore, err := store.NewStore(directoryPath, path.Join(directoryPath, "import-backup"), 0)
	if err != nil {
		return dataImport.Result{}, err
	}

	ctx.Store.StartBatch()
	defer ctx.Store.StopBatch()

	createdFrames := 0
	createdTags := 0
	createdProjects := 0
	reusedProjects := 0

	topLevelProject, err := ctx.Store.AddProject(model.Project{
		Name: fmt.Sprintf("Import %s", time.Now().Format(time.RFC3339)),
	})

	// create projects, don't match by name
	projectMapping := make(map[string]string)
	for _, p := range importStore.Projects() {
		if created, reused, err := importProject(*p, importStore, ctx, projectMapping, topLevelProject); err != nil {
			return dataImport.Result{}, err
		} else {
			createdProjects += created
			reusedProjects += reused
		}
	}

	// import tags, match by name
	for _, t := range importStore.Tags() {
		existingTag, _ := ctx.Store.FindFirstTag(func(tag *model.Tag) bool {
			return tag.Name == t.Name
		})

		if existingTag == nil {
			if _, err := ctx.Store.AddTag(*t); err != nil {
				return dataImport.Result{}, err
			}
			createdTags += 1
		}
	}

	// import frames
	for _, f := range importStore.Frames() {
		created, err := importFrame(*f, projectMapping, importStore, ctx)
		if err != nil {
			return dataImport.Result{}, err
		}
		createdFrames += created
	}

	return dataImport.Result{
		CreatedProjects: createdProjects,
		ReusedProjects:  reusedProjects,
		CreatedTags:     createdTags,
		CreatedFrames:   createdFrames,
	}, err
}

func importProject(p model.Project, importStore model.Store, ctx *context.TomContext, mapping map[string]string, topLevelProject *model.Project) (int, int, error) {
	createdProjects := 0
	reusedProjects := 0

	// already imported
	if mapping[p.ID] != "" {
		return 0, 1, nil
	}

	if parent := p.Parent(); parent != nil {
		created, reused, err := importProject(*parent, importStore, ctx, mapping, topLevelProject)
		if err != nil {
			return createdProjects, reusedProjects, err
		}

		p.ParentID = mapping[p.ParentID]

		createdProjects += created
		reusedProjects += reused
	} else {
		p.ParentID = topLevelProject.ID
	}

	imported, err := ctx.Store.AddProject(p)
	mapping[p.ID] = imported.ID
	return 1, 0, err
}

func importFrame(f model.Frame, projectMapping map[string]string, importStore model.Store, ctx *context.TomContext) (int, error) {
	mappedProject := projectMapping[f.ProjectId]
	if mappedProject == "" {
		p, _ := importStore.ProjectByID(f.ProjectId)
		return 0, fmt.Errorf("unable to find project for frame: %s, %s, %v\n", f.ID, f.ProjectId, p)
	}

	f.ProjectId = mappedProject
	if _, err := ctx.Store.AddFrame(f); err != nil {
		return 0, err
	}
	return 1, nil
}
