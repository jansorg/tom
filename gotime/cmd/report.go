package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"../context"
	"../htmlreport"
	"../report"
)

func newReportCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var fromDateString string
	var toDateString string

	var day int8
	var month int8
	var year int

	var roundFrames time.Duration
	var roundNearest bool
	var roundTotal time.Duration
	var roundTotalNearest bool

	var htmlFile string

	var cmd = &cobra.Command{
		Use:   "report",
		Short: "Reporting about the tracked time",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var start time.Time
			var end time.Time

			if fromDateString != "" {
				start, err = time.Parse(time.RFC3339Nano, fromDateString)
				if err != nil {
					fatal(err)
				}
			}

			if toDateString != "" {
				end, err = time.Parse(time.RFC3339Nano, toDateString)
				if err != nil {
					fatal(err)
				}
			}

			if cmd.Flag("day") != nil {
				now := time.Now()
				year, month, today := now.Date()
				start = time.Date(year, month, today+int(day), 0, 0, 0, 0, now.Location())
				end = time.Date(year, month, today+int(day), 24, 0, 0, 0, now.Location())
			} else if cmd.Flag("month") != nil {
				now := time.Now()
				year, currentMonth, _ := now.Date()
				start = time.Date(year, time.Month(int(currentMonth)+int(month)), 0, 0, 0, 0, 0, now.Location())
				end = time.Date(year, time.Month(int(currentMonth)+int(month)+1), 0, 0, 0, 0, 0, now.Location())
			} else if cmd.Flag("year") != nil {
				now := time.Now()
				currentYear, _, _ := now.Date()
				start = time.Date(currentYear+year, time.January, 0, 0, 0, 0, 0, now.Location())
				end = time.Date(currentYear+year, time.December, 24, 0, 0, 0, 0, now.Location())
			}

			var frameRoundingMode = report.RoundNone
			if roundFrames > 0 {
				if roundNearest {
					frameRoundingMode = report.RoundNearest
				} else {
					frameRoundingMode = report.RoundUp
				}
			}

			var roundingModeTotal = report.RoundNone
			if roundTotal > 0 {
				if roundTotalNearest {
					roundingModeTotal = report.RoundNearest
				} else {
					roundingModeTotal = report.RoundUp
				}
			}

			frameReport := report.TimeReport{
				FrameRoundingMode: frameRoundingMode,
				RoundFramesTo:     roundFrames,
				TotalRoundingMode: roundingModeTotal,
				RoundTotalTo:      roundTotal,
			}

			var usedStart *time.Time
			if !start.IsZero() {
				usedStart = &start
			}

			var usedEnd *time.Time
			if !end.IsZero() {
				usedEnd = &end
			}

			result, err := frameReport.Calc(usedStart, usedEnd, context)
			if err != nil {
				fatal(err)
			}

			if context.JsonOutput {
				data, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					fatal(err)
				}
				fmt.Println(string(data))
			} else if htmlFile != "" {
				templatePath := filepath.Join("../templates", "reports/default.gohtml")
				if templatePath, err = filepath.Abs(templatePath); err != nil {
					fatal(err)
				}
				htmlReport := htmlreport.NewReport(templatePath)

				content, err := htmlReport.Render(result)
				if err != nil {
					fatal(err)
				}

				err = ioutil.WriteFile(htmlFile, []byte(content), 0600)
				if err != nil {
					fatal(err)
				}
			} else {
				if result.From != nil {
					fmt.Printf("From: %s\n", result.From.String())
				}

				if result.To != nil {
					fmt.Printf("To: %s\n", result.To.String())
				}

				for _, r := range result.Items {
					fmt.Printf("%s: %s\n", r.Name, r.Duration.String())
				}
			}
		},
	}

	cmd.Flags().StringVarP(&fromDateString, "from", "f", "", "Optional start date")
	cmd.Flags().StringVarP(&toDateString, "to", "t", "", "Optional end date")
	cmd.Flags().Int8VarP(&day, "day", "", 0, "Select the date range of a given day. For example, 0 is today, -1 is one day ago, etc.")
	cmd.Flags().Int8VarP(&month, "month", "", 0, "Filter on a given month. For example, 0 is the current month, -1 is last month, etc.")
	cmd.Flags().IntVarP(&year, "year", "", 0, "Filter on a specific year. 0 is the current year, -1 is last year, etc.")

	cmd.Flags().DurationVarP(&roundFrames, "round-frames", "r", time.Duration(0), "Round durations of each frame to the nearest multiple of this duration")
	cmd.Flags().BoolVarP(&roundNearest, "round-frames-nearest", "", false, "Round the durations of each frame to the nearest multiple. The default is to round up.")

	cmd.Flags().DurationVarP(&roundTotal, "round-total", "", time.Duration(0), "Round the overall duration of each project to the next matching multiple of this duration")
	cmd.Flags().BoolVarP(&roundTotalNearest, "round-total-nearest", "", false, "Round the overall duration of a project or tag the nearest multiple. The default is to round up.")

	cmd.Flags().StringVarP(&htmlFile, "html", "", "", "Output the report as HTML into the given file")

	parent.AddCommand(cmd)
	return cmd
}
