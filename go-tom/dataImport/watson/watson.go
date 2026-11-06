package watson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dataImport"
	"github.com/jansorg/tom/go-tom/model"
)

var epoch = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

func NewImporter() dataImport.Handler {
	return &watsonImporter{}
}

type watsonImporter struct{}

func (watsonImporter) Import(filename string, ctx *context.TomContext) (result dataImport.Result, err error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return dataImport.Result{}, err
	}

	var frames [][]interface{}
	if err := json.Unmarshal(bytes, &frames); err != nil {
		return dataImport.Result{}, fmt.Errorf("error parsing json file %s: %s", filename, err.Error())
	}

	createdProjects := 0
	reusedProjects := 0
	createdTags := 0
	createdFrames := 0

	ctx.Store.StartBatch()
	defer func() {
		if stopErr := ctx.Store.StopBatch(); err == nil && stopErr != nil {
			err = stopErr
		}
	}()

	for i, f := range frames {
		start, stop, projectName, tagNames, updated, err := parseFrame(f)
		if err != nil {
			return dataImport.Result{}, fmt.Errorf("invalid frame #%d in %s: %v", i+1, filename, err)
		}

		startTime := epoch.Add(time.Duration(start) * time.Second)
		stopTime := epoch.Add(time.Duration(stop) * time.Second)
		udpatedTime := epoch.Add(time.Duration(updated) * time.Second)

		project, createdProject, err := ctx.StoreHelper.GetOrCreateNestedProject(projectName)
		if err != nil {
			return dataImport.Result{}, err
		}

		if createdProject {
			createdProjects++
		} else {
			reusedProjects++
		}

		var tagIDs []string
		for _, tagName := range tagNames {
			tag, createdTag, err := ctx.StoreHelper.GetOrCreateTag(tagName)
			if err != nil {
				return dataImport.Result{}, err
			}
			tagIDs = append(tagIDs, tag.ID)

			if createdTag {
				createdTags++
			}
		}

		_, err = ctx.Store.AddFrame(model.Frame{
			Start:     &startTime,
			End:       &stopTime,
			Updated:   &udpatedTime,
			ProjectId: project.ID,
			TagIDs:    tagIDs,
		})
		if err != nil {
			return dataImport.Result{}, err
		}
		createdFrames++
	}

	return dataImport.Result{
		CreatedProjects: createdProjects,
		CreatedTags:     createdTags,
		CreatedFrames:   createdFrames,
	}, err
}

// parseFrame parses a Watson frame: [start, stop, project, id, tags, updated]
// tags may be null, updated is optional and defaults to stop
func parseFrame(f []interface{}) (start, stop float64, project string, tags []string, updated float64, err error) {
	if len(f) < 5 {
		return 0, 0, "", nil, 0, fmt.Errorf("expected at least 5 values, found %d", len(f))
	}

	var ok bool
	if start, ok = f[0].(float64); !ok {
		return 0, 0, "", nil, 0, fmt.Errorf("start is not a number: %v", f[0])
	}
	if stop, ok = f[1].(float64); !ok {
		return 0, 0, "", nil, 0, fmt.Errorf("stop is not a number: %v", f[1])
	}
	if project, ok = f[2].(string); !ok || project == "" {
		return 0, 0, "", nil, 0, fmt.Errorf("project is not a valid name: %v", f[2])
	}

	if f[4] != nil {
		tagValues, ok := f[4].([]interface{})
		if !ok {
			return 0, 0, "", nil, 0, fmt.Errorf("tags are not a list: %v", f[4])
		}
		for _, v := range tagValues {
			name, ok := v.(string)
			if !ok {
				return 0, 0, "", nil, 0, fmt.Errorf("tag is not a string: %v", v)
			}
			tags = append(tags, name)
		}
	}

	updated = stop
	if len(f) >= 6 && f[5] != nil {
		if updated, ok = f[5].(float64); !ok {
			return 0, 0, "", nil, 0, fmt.Errorf("updated is not a number: %v", f[5])
		}
	}

	return start, stop, project, tags, updated, nil
}
