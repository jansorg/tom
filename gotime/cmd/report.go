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
	"github.com/jansorg/gotime/gotime/htmlreport"
	"github.com/jansorg/gotime/gotime/report"
)

func newReportCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var fromDateString string
	var toDateString string
	var projectID string

	var day int
	var month int
	var year int

	var splitModes string

	var roundFrames time.Duration
	var roundTotals time.Duration

	var roundModeFrames string
	var roundModeTotal string

	var templatePath string

	var cmd = &cobra.Command{
		Use:   "report",
		Short: "Generate reports about your tracked time",
		Run: func(cmd *cobra.Command, args []string) {
			filterRange := dateUtil.NewDateRange(nil, nil, ctx.Locale)

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
				filterRange = dateUtil.NewDayRange(time.Now(), ctx.Locale).Shift(0, 0, day)
			} else if cmd.Flag("month").Changed {
				filterRange = dateUtil.NewMonthRange(time.Now(), ctx.Locale).Shift(0, month, 0)
			} else if cmd.Flag("year").Changed {
				filterRange = dateUtil.NewYearRange(time.Now(), ctx.Locale).Shift(year, 0, 0)
			}

			var frameRoundingMode = dateUtil.ParseRoundingMode(roundModeFrames)
			var totalsRoundingNode = dateUtil.ParseRoundingMode(roundModeTotal)

			var splitOperations []report.SplitOperation
			if splitModes != "" {
				for _, mode := range strings.Split(splitModes, ",") {
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
			}

			// project filter
			if projectID != "" {
				// if it's a name resolve it to the ID
				if project, err := context.Query.ProjectByFullName(projectID); err == nil {
					projectID = project.ID
				}
				if _, err := context.Query.ProjectByID(projectID); err != nil {
					fatal(fmt.Errorf("project %s not found", projectID))
				}
			}

			storeFrames := context.Store.Frames()
			frameReport := report.NewBucketReport(frames.NewSortedFrameList(storeFrames), context)
			frameReport.ProjectID = projectID
			frameReport.FilterRange = filterRange
			frameReport.RoundFramesTo = roundFrames
			frameReport.RoundTotalsTo = roundTotals
			frameReport.RoundingModeFrames = frameRoundingMode
			frameReport.RoundingModeTotals = totalsRoundingNode
			frameReport.SplitOperations = splitOperations
			frameReport.Update()

			if templatePath != "" {
				if err := printTemplate(context, templatePath, frameReport); err != nil {
					fatal(fmt.Errorf("error rendering with template: %s", err.Error()))
				}
			} else if context.JsonOutput {
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

	templateAnnotations := make(map[string][]string)
	templateAnnotations[cobra.BashCompFilenameExt] = []string{"gohtml"}
	cmd.Flags().StringVarP(&templatePath, "template", "", "", "Template to use for rendering. This may either be a full path to a template file or the name (without extension) of a template shipped with gotime.")
	cmd.Flag("template").Annotations = templateAnnotations

	cmd.Flags().StringVarP(&fromDateString, "from", "f", "", "Optional start date")
	cmd.Flags().StringVarP(&toDateString, "to", "t", "", "Optional end date")
	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project filter. Only frames which belong to this project are used for the report.")

	cmd.Flags().IntVarP(&day, "day", "", 0, "Select the date range of a given day. For example, 0 is today, -1 is one day ago, etc.")
	cmd.Flags().IntVarP(&month, "month", "", 0, "Filter on a given month. For example, 0 is the current month, -1 is last month, etc.")
	cmd.Flags().IntVarP(&year, "year", "", 0, "Filter on a specific year. 0 is the current year, -1 is last year, etc.")

	cmd.Flags().StringVarP(&splitModes, "split", "", "", "Group frames into years, months and/or days. Possible values: year,month,day")

	cmd.Flags().DurationVarP(&roundFrames, "round-frames-to", "", time.Duration(0), "Round durations of each frame to the nearest multiple of this duration")
	cmd.Flags().StringVarP(&roundModeFrames, "round-frames", "", "up", "Rounding mode for sums of durations. Default: up. Possible values: up|nearest")

	// fixme
	cmd.Flags().DurationVarP(&roundTotals, "round-totals-to", "", time.Duration(0), "Round the overall duration of each project to the next matching multiple of this duration")
	cmd.Flags().StringVarP(&roundModeTotal, "round-totals", "", "up", "Rounding mode for sums of durations. Default: up. Possible values: up|nearest")

	parent.AddCommand(cmd)
	return cmd
}

func printTemplate(ctx *context.GoTimeContext, templatePath string, report *report.BucketReport) error {
	t := htmlreport.NewReport(templatePath, ctx)
	out, err := t.Render(report)
	if err != nil {
		return err
	}

	fmt.Println(out)
	return nil
}

func parseDate(dateString *string) (*time.Time, error) {
	result, err := time.Parse(time.RFC3339Nano, *dateString)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func printReport(report *report.ResultBucket, ctx *context.GoTimeContext, level int) {
	title := report.Title()
	if title != "" {
		printlnIndenting(level-1, title)
	} else if level == 1 {
		printlnIndenting(level-1, "Overall")
	}

	if !report.DateRange.Empty() {
		printfIndenting(level, "Date range: %s\n", report.DateRange.ShortString())
	}
	if !report.TrackedDateRange.Empty() {
		printfIndenting(level, "Tracked time: %s\n", report.TrackedDateRange.ShortString())
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
