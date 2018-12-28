package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/store"
)

func newImportFanurioCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "fanurio",
		Short: "import frames and projects from Fanurio CSV output",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			absPath, err := filepath.Abs(args[0])
			if err != nil {
				fatal(err)
			}

			if err = importCSV(absPath, ctx); err != nil {
				fatal(err)
			}
		},
	}

	parent.AddCommand(cmd)
	return cmd
}

func importCSV(filePath string, ctx *context.GoTimeContext) error {
	query := ctx.Query
	dataStore := ctx.Store

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

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
			return err
		}

		endTime, err := parseTime(fmt.Sprintf("%s %s", dateString, endString))
		if err != nil {
			return err
		}

		fullProjectName := strings.Join([]string{clientName, projectName, taskName}, "/")
		project, err := query.ProjectByFullName(fullProjectName)
		if err != nil {
			project, _ = dataStore.AddProject(store.Project{Name: fullProjectName})
		}

		_, err = dataStore.AddFrame(store.Frame{
			ProjectId: project.ID,
			Notes:     notes,
			Start:     &startTime,
			End:       &endTime,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func logError(err error) {
	log.Printf("Import error: %v", err)
}

func parseTime(value string) (time.Time, error) {
	d, err := date.ParseTime("02.01.06 15:04:05", value)
	if err != nil {
		d, err = date.ParseTime("02.01.06 15:04", value)
	}
	return d, err
}
