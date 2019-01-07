package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/go-tom/context"
	"github.com/jansorg/gotime/go-tom/store"
)

func newImportWatsonCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "watson",
		Short: "",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dir := ""
			if len(args) == 1 {
				dir = args[1]
			} else {
				home, err := homedir.Dir()
				if err != nil {
					fatal(err)
				}
				dir = filepath.Join(home, ".config", "watson")
			}

			framesFile := filepath.Join(dir, "frames")
			if _, err := os.Stat(framesFile); os.IsNotExist(err) {
				fatal(fmt.Errorf("file %s not found", framesFile))
			}

			var frames [][]interface{}
			bytes, err := ioutil.ReadFile(framesFile)
			if err != nil {
				fatal(err)
			}
			if err := json.Unmarshal(bytes, &frames); err != nil {
				fatal(fmt.Errorf("error parsing json file %s: %s", framesFile, err.Error()))
			}

			epoch := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

			createdProjects := 0
			createdTags := 0

			ctx.Store.StartBatch()
			defer ctx.Store.StopBatch()

			for _, f := range frames {
				start := f[0].(float64)
				stop := f[1].(float64)
				projectName := f[2].(string)
				// frameID := f[3].(string)
				tagNames := f[4].([]interface{})
				updated := f[5].(float64)

				startTime := epoch.Add(time.Duration(start) * time.Second)
				stopTime := epoch.Add(time.Duration(stop) * time.Second)
				udpatedTime := epoch.Add(time.Duration(updated) * time.Second)

				project, createdProject, err := ctx.StoreHelper.GetOrCreateNestedProject(projectName)
				if err != nil {
					fatal(err)
				}

				if createdProject {
					createdProjects++
				}

				var tagIDs []string
				for _, tagName := range tagNames {
					tag, createdTag, err := ctx.StoreHelper.GetOrCreateTag(tagName.(string))
					if err != nil {
						fatal(err)
					}
					tagIDs = append(tagIDs, tag.ID)

					if createdTag {
						createdTags++
					}
				}

				_, err = ctx.Store.AddFrame(store.Frame{
					Start:     &startTime,
					End:       &stopTime,
					Updated:   &udpatedTime,
					ProjectId: project.ID,
					TagIDs:    tagIDs,
				})
				if err != nil {
					fatal(err)
				}
			}

			fmt.Printf("Successfully imported %d projects, %d tags, and %d frames.\n", createdProjects, createdTags, len(frames))
		},
	}

	parent.AddCommand(cmd)
	return cmd
}
