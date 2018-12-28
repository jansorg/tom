package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

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

		projectName := row[0]
		notes := row[1]
		dateString := row[2]
		startString := row[3]
		endString := row[4]

		// Mon Jan 2 15:04:05 MST 2006
		startTime, err := date.ParseTime("02.01.06 15:04:05", fmt.Sprintf("%s %s", dateString, startString))
		if err != nil {
			return err
		}

		endTime, err := date.ParseTime("02.01.06 15:04:05", fmt.Sprintf("%s %s", dateString, endString))
		if err != nil {
			return err
		}

		project, err := query.ProjectByFullName(projectName)
		if err != nil {
			project, _ = dataStore.AddProject(store.Project{ShortName: projectName, FullName: projectName})
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
