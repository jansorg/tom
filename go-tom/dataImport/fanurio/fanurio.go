package fanurio

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dataImport"
	"github.com/jansorg/tom/go-tom/model"
)

func NewCSVImporter() dataImport.Handler {
	return &fanurioImporter{}
}

type fanurioImporter struct{}

func (fanurioImporter) Import(filePath string, ctx *context.TomContext) (dataImport.Result, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return dataImport.Result{}, err
	}
	defer file.Close()

	ctx.Store.StartBatch()
	defer ctx.Store.StopBatch()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	rows, err := reader.ReadAll()
	if err != nil {
		return dataImport.Result{}, err
	}

	createdFrames := 0
	createdProjects := 0
	reusedProjects := 0

	for i, row := range rows {
		if i == 0 {
			continue
		}

		clientName := strings.TrimSpace(row[0])
		projectName := strings.TrimSpace(row[1])
		taskName := strings.TrimSpace(row[2])
		notes := strings.TrimSpace(row[3])
		dateString := row[4]
		startString := row[5]
		endString := row[6]

		// Mon Jan 2 15:04:05 MST 2006
		startTime, err := parseTime(fmt.Sprintf("%s %s", dateString, startString))
		if err != nil {
			return dataImport.Result{}, err
		}

		endTime, err := parseTime(fmt.Sprintf("%s %s", dateString, endString))
		if err != nil {
			return dataImport.Result{}, err
		}

		project, created, err := ctx.StoreHelper.GetOrCreateNestedProjectNames(clientName, projectName, taskName)
		if err != nil {
			return dataImport.Result{}, err
		}

		if created {
			createdProjects++
		} else {
			reusedProjects++
		}

		_, err = ctx.Store.AddFrame(model.Frame{
			ProjectId: project.ID,
			Notes:     notes,
			Start:     &startTime,
			End:       &endTime,
		})
		if err != nil {
			return dataImport.Result{}, err
		}

		createdFrames++
	}

	return dataImport.Result{
		CreatedProjects: createdProjects,
		ReusedProjects:  reusedProjects,
		CreatedFrames:   createdFrames,
	}, err
}

func parseTime(value string) (time.Time, error) {
	d, err := time.Parse("02.01.06 15:04:05", value)
	if err != nil {
		d, err = time.Parse("02.01.06 15:04", value)
	}
	return d, err
}
