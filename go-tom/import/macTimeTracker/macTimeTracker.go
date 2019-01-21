package macTimeTracker

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

func ImportCSV(filename string, ctx *context.GoTimeContext) (created int, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	ctx.Store.StartBatch()
	defer ctx.Store.StopBatch()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	reader.Comma = ';'

	rows, err := reader.ReadAll()
	if err != nil {
		return 0, err
	}

	createdFrames := 0
	for i, row := range rows {
		if i == 0 {
			// ignore header
			continue
		}

		if len(row) != 7 {
			return 0, fmt.Errorf("unexpected number of columns %d instead of 7 expected", len(row))
		}

		projectName := strings.TrimSpace(row[0])
		taskName := strings.TrimSpace(row[1])
		// ignore start date
		startString := row[3] // 2015-12-18 17:00:37
		// ignore end date, using start+duration instead
		// 00:23:17, the duration is sometimes different from end-start (off by 1s), we assume that TimeTracker is tracking in ms values and is rounding the duration
		durationString := row[5]
		// ignoring duration column
		notes := strings.TrimSpace(row[6])

		startTime, err := parseTime(startString)
		if err != nil {
			return 0, err
		}

		duration, err := parseDuration(durationString)
		if err != nil {
			return 0, err
		}

		endTime := startTime.Add(duration)

		project, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames(projectName, taskName)
		if err != nil {
			return 0, err
		}

		if _, err = ctx.Store.AddFrame(model.Frame{
			ProjectId: project.ID,
			Notes:     notes,
			Start:     &startTime,
			End:       &endTime,
		}); err != nil {
			return 0, err
		}

		createdFrames++
	}

	return createdFrames, nil
}

func parseTime(value string) (time.Time, error) {
	d, err := date.ParseTime("2006-01-02 15:04:05", value)
	return d, err
}

func parseDuration(value string) (time.Duration, error) {
	values := strings.Split(value, ":")

	var hours, minutes, seconds int
	var err error
	if hours, err = strconv.Atoi(values[0]); err != nil {
		return 0, err
	}
	if minutes, err = strconv.Atoi(values[1]); err != nil {
		return 0, err
	}
	if seconds, err = strconv.Atoi(values[2]); err != nil {
		return 0, err
	}

	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}
