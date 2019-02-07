package cmdUtil

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/util"
)

type PropList interface {
	Size() int
	Get(index int, prop string, format string, ctx* context.TomContext) (interface{}, error)
}

func AddListOutputFlags(cmd *cobra.Command, defaultFormat string, supportedProps []string) {
	cmd.Flags().StringP("output", "o", "plain", "Output format. Supported: plain | json. Default: plain")
	cmd.Flags().StringP("format", "f", defaultFormat, fmt.Sprintf("A comma separated list of of properties to output. Default: %s. Possible values: %s", defaultFormat, strings.Join(supportedProps, ",")))
	cmd.Flags().StringP("delimiter", "d", "\t", "The delimiter to add between property values. Default: TAB")
}

func parseListOutputFlags(cmd *cobra.Command) (props []string, output string, delimiter string, err error) {
	output, err = cmd.Flags().GetString("output")
	if err != nil {
		return nil, "", "", err
	}

	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, "", "", err
	}

	props = strings.Split(format, ",")
	delimiter, err = cmd.Flags().GetString("delimiter")
	if err != nil {
		return nil, "", "", err
	}

	return props, output, delimiter, nil
}

func PrintList(cmd *cobra.Command, data PropList, ctx *context.TomContext) error {
	formatFlags, output, delimiter, err := parseListOutputFlags(cmd)
	if err != nil {
		util.Fatal(err)
	}

	type row map[string]interface{}

	var rows []row
	for i := 0; i < data.Size(); i++ {
		r := row{}
		for _, prop := range formatFlags {
			r[prop], err = data.Get(i, prop, output, ctx)
			if err != nil {
				return err
			}
		}
		rows = append(rows, r)
	}

	switch output {
	case "plain":
		for _, row := range rows {
			var rowValues []string
			for _, prop := range formatFlags {
				rowValues = append(rowValues, stringValue(row[prop], ctx))
			}
			fmt.Println(strings.Join(rowValues, delimiter))
		}
	case "json":
		PrintJSON(rows)
	default:
		util.Fatal(fmt.Errorf("unsupported output type %s", output))
	}
	return nil
}

func PrintJSON(value interface{}) {
	if bytes, err := json.MarshalIndent(value, "", "  "); err != nil {
		util.Fatal(err)
	} else {
		fmt.Println(string(bytes))
	}
}

func stringValue(v interface{}, ctx *context.TomContext) string {
	if s, ok := v.(string); ok {
		return s
	}

	if date, ok := v.(time.Time); ok {
		return date.Format(time.RFC3339)
	}

	if duration, ok := v.(time.Duration); ok {
		return strconv.FormatInt(duration.Nanoseconds()/1000/1000, 10)
	}

	if date, ok := v.(util.DateRange); ok {
		return date.ShortString()
	}

	return fmt.Sprintf("%v", v)
}
