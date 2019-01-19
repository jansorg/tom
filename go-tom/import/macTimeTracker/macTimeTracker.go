package macTimeTracker

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"howett.net/plist"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

func Import(filename string, ctx *context.GoTimeContext) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	decoder := plist.NewDecoder(file)
	data := map[string]interface{}{}
	err = decoder.Decode(&data)

	if err == nil {
		for _, v := range data {
			// fmt.Printf("%v\n", k)

			if m, ok := v.([]interface{}); ok {
				data := m[18]
				if d, ok := data.(map[string]interface{}); ok {
					entries, ok := d["NS.data"]
					if ok {
						if b, ok := entries.([]byte); ok {
							dec := plist.NewDecoder(bytes.NewReader(b))
							list := map[string]interface{}{}
							err = dec.Decode(&list)
							if err != nil {
								fmt.Println("error reading data bytes")
							} else {
								for k1, v1 := range list {
									fmt.Printf("%s = %v\n\n", k1, v1)
								}
							}
						}
					}
				}
			}
		}
	} else {
		fmt.Printf("error: %v", err)
	}

	return err
}

func ImportCSV(filename string, ctx *context.GoTimeContext) (created int, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}

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
		endString := row[4]   // 2015-12-18 17:00:37
		// ignoring duration column
		notes := strings.TrimSpace(row[6])

		startTime, err := parseTime(startString)
		if err != nil {
			return 0, err
		}

		endTime, err := parseTime(endString)
		if err != nil {
			return 0, err
		}

		project, _, err := ctx.StoreHelper.GetOrCreateNestedProjectNames(projectName, taskName)
		if err != nil {
			return 0, err
		}

		_, err = ctx.Store.AddFrame(model.Frame{
			ProjectId: project.ID,
			Notes:     notes,
			Start:     &startTime,
			End:       &endTime,
		})
		if err != nil {
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
