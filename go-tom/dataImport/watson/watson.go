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

func (watsonImporter) Import(filename string, ctx *context.GoTimeContext) (dataImport.Result, error) {
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
	defer ctx.Store.StopBatch()

	for _, f := range frames {
		start := f[0].(float64)
		stop := f[1].(float64)
		projectName := f[2].(string)
		// ignoring frameID f[3]
		tagNames := f[4].([]interface{})
		updated := f[5].(float64)

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
			tag, createdTag, err := ctx.StoreHelper.GetOrCreateTag(tagName.(string))
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
		createdFrames++;
	}

	return dataImport.Result{
		CreatedProjects: createdProjects,
		CreatedTags:     createdTags,
		CreatedFrames:   createdFrames,
	}, err
}
