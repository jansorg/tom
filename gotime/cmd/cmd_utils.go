package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func addListOutputFlags(cmd *cobra.Command, supportedProps []string) {
	cmd.Flags().StringP("output", "o", "plain", "Output format. Supported: plain | json. Default: plain")
	cmd.Flags().StringP("format", "f", "name", fmt.Sprintf("A comma separated list of of properties to output. Default: id . Possible values: %s", strings.Join(supportedProps, ",")))
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

type propList interface {
	size() int
	get(index int, prop string) (string, error)
}

func printList(cmd *cobra.Command, data propList) error {
	formatFlags, output, delimiter, err := parseListOutputFlags(cmd)
	if err != nil {
		fatal(err)
	}

	type row map[string]string

	var rows []row
	for i := 0; i < data.size(); i++ {
		r := row{}
		for _, prop := range formatFlags {
			r[prop], err = data.get(i, prop)
			if err != nil {
				return err
			}
		}
		rows = append(rows, r)
	}

	switch output {
	case "plain":
		for _, row := range rows {
			rowValues := []string{}
			for _, prop := range formatFlags {
				rowValues = append(rowValues, row[prop])
			}
			fmt.Println(strings.Join(rowValues, delimiter))
		}
	case "json":
		if bytes, err := json.MarshalIndent(rows, "", "  "); err != nil {
			fatal(err)
		} else {
			fmt.Println(string(bytes))
		}
	default:
		fatal(fmt.Errorf("unsupported output type %s", output))
	}
	return nil
}
