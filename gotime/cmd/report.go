package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/dateUtil"
	"github.com/jansorg/gotime/gotime/frames"
	"github.com/jansorg/gotime/gotime/report"
)

func newReportCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var fromDateString string
	var toDateString string

	var day int
	var month int
	var year int

	var splitModes []string

	var roundFrames time.Duration
	var roundTotals time.Duration

	var roundModeFrames string
	var roundModeTotal string

	var cmd = &cobra.Command{
		Use:   "report",
		Short: "Generate reports about your tracked time",
		Run: func(cmd *cobra.Command, args []string) {
			var filterRange dateUtil.DateRange

			if fromDateString != "" {
				start, err := parseDate(&fromDateString)
				if err != nil {
					fatal(err)
				}
				filterRange.Start = start
			}

			if toDateString != "" {
				end, err := parseDate(&toDateString)
				if err != nil {
					fatal(err)
				}
				filterRange.End = end
			}

			// day, month, year params override the filter values
			if cmd.Flag("day").Changed {
				filterRange = dateUtil.NewDayRange(time.Now()).Shift(0, 0, day)
			} else if cmd.Flag("month").Changed {
				filterRange = dateUtil.NewMonthRange(time.Now()).Shift(0, month, 0)
			} else if cmd.Flag("year").Changed {
				filterRange = dateUtil.NewYearRange(time.Now()).Shift(year, 0, 0)
			}

			var frameRoundingMode = dateUtil.ParseRoundingMode(roundModeFrames)
			var totalsRoundingNode = dateUtil.ParseRoundingMode(roundModeTotal)

			var splitOperations []report.SplitOperation
			for _, mode := range splitModes {
				switch mode {
				case "year":
					splitOperations = append(splitOperations, report.SplitByYear)
				case "month":
					splitOperations = append(splitOperations, report.SplitByMonth)
				case "day":
					splitOperations = append(splitOperations, report.SplitByDay)
				case "project":
					splitOperations = append(splitOperations, report.SplitByProject)
				default:
					fatal(fmt.Errorf("unknown split value %s. Supported: year, month, day", mode))
				}
			}

			storeFrames := context.Store.Frames()
			frameReport := report.NewBucketReport(frames.NewFrameList(storeFrames))
			frameReport.FilterRange = filterRange
			frameReport.RoundFramesTo = roundFrames
			frameReport.RoundTotalsTo = roundTotals
			frameReport.RoundingModeFrames = frameRoundingMode
			frameReport.RoundingModeTotals = totalsRoundingNode
			frameReport.SplitOperations = splitOperations
			frameReport.Update()

			if context.JsonOutput {
				data, err := json.MarshalIndent(frameReport.Result, "", "  ")
				if err != nil {
					fatal(err)
				}
				fmt.Println(string(data))
			} else {
				printReport(frameReport.Result, context, 1)
			}
		},
	}

	cmd.PersistentFlags().StringVarP(&fromDateString, "from", "f", "", "Optional start date")
	cmd.PersistentFlags().StringVarP(&toDateString, "to", "t", "", "Optional end date")
	cmd.PersistentFlags().IntVarP(&day, "day", "", 0, "Select the date range of a given day. For example, 0 is today, -1 is one day ago, etc.")
	cmd.PersistentFlags().IntVarP(&month, "month", "", 0, "Filter on a given month. For example, 0 is the current month, -1 is last month, etc.")
	cmd.PersistentFlags().IntVarP(&year, "year", "", 0, "Filter on a specific year. 0 is the current year, -1 is last year, etc.")

	cmd.PersistentFlags().StringArrayVarP(&splitModes, "group", "", []string{}, "Group frames into years, months and/or days. Possible values: year,month,day")

	cmd.PersistentFlags().DurationVarP(&roundFrames, "round-frames-to", "", time.Duration(0), "Round durations of each frame to the nearest multiple of this duration")
	cmd.PersistentFlags().StringVarP(&roundModeFrames, "round-frames", "", "up", "Rounding mode for sums of durations. Default: up. Possible values: up|nearest")

	cmd.PersistentFlags().DurationVarP(&roundTotals, "round-totals-to", "", time.Duration(0), "Round the overall duration of each project to the next matching multiple of this duration")
	cmd.PersistentFlags().StringVarP(&roundModeTotal, "round-totals", "", "up", "Rounding mode for sums of durations. Default: up. Possible values: up|nearest")

	parent.AddCommand(cmd)
	newReportHtmlCommand(context, cmd)

	return cmd
}

func parseDate(dateString *string) (*time.Time, error) {
	result, err := time.Parse(time.RFC3339Nano, *dateString)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func printReport(report *report.ResultBucket, ctx *context.GoTimeContext, level int) {
	title := report.Title(ctx)
	if title != "" {
		printlnIndenting(level-1, title)
	} else if level == 1 {
		printlnIndenting(level-1, "Overall")
	}

	if !report.DateRange.Empty() {
		printlnIndenting(level, report.DateRange.String())
	}
	if !report.UsedDateRange.Empty() {
		printlnIndenting(level, report.UsedDateRange.String())
	}
	printfIndenting(level, "Duration: %s\n", report.Duration.String())
	printfIndenting(level, "Exact Duration: %s\n", report.ExactDuration.String())
	fmt.Println()

	for _, r := range report.Results {
		printReport(r, ctx, level+1)
	}
}

func printlnIndenting(level int, value string) {
	fmt.Print(strings.Repeat("    ", level))
	fmt.Println(value)
}

func printfIndenting(level int, format string, a ...interface{}) {
	fmt.Print(strings.Repeat("    ", level))
	fmt.Printf(format, a...)
}
