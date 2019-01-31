package report

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateUtil"
	"github.com/jansorg/tom/go-tom/htmlreport"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/report"
)

func NewCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var includeActiveFrames bool

	var jsonOutput bool

	var fromDateString string
	var toDateString string
	var projectFilter []string

	var day int
	var month int
	var year int

	var splitModes string

	var roundFrames time.Duration
	var roundTotals time.Duration

	var roundModeFrames string
	var roundModeTotal string

	var decimalDurations bool

	var templateName string
	var templateFilePath string

	var cmd = &cobra.Command{
		Use:   "report",
		Short: "Generate reports about your tracked time",
		Run: func(cmd *cobra.Command, args []string) {
			filterRange := dateUtil.NewDateRange(nil, nil, ctx.Locale)

			if fromDateString != "" {
				start, err := parseDate(&fromDateString)
				if err != nil {
					log.Fatal(err)
				}
				filterRange.Start = start
			}

			if toDateString != "" {
				end, err := parseDate(&toDateString)
				if err != nil {
					log.Fatal(err)
				}
				filterRange.End = end
			}

			// day, month, year params override the filter values
			if cmd.Flag("day").Changed {
				filterRange = dateUtil.NewDayRange(time.Now(), ctx.Locale, time.Local).Shift(0, 0, day)
			} else if cmd.Flag("month").Changed {
				filterRange = dateUtil.NewMonthRange(time.Now(), ctx.Locale, time.Local).Shift(0, month, 0)
			} else if cmd.Flag("year").Changed {
				filterRange = dateUtil.NewYearRange(time.Now(), ctx.Locale, time.Local).Shift(year, 0, 0)
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
					case "week":
						splitOperations = append(splitOperations, report.SplitByWeek)
					case "day":
						splitOperations = append(splitOperations, report.SplitByDay)
					case "project":
						splitOperations = append(splitOperations, report.SplitByProject)
					// case "parentProject":
					// 	splitOperations = append(splitOperations, report.SplitByParentProject)
					default:
						log.Fatal(fmt.Errorf("unknown split value %s. Supported: year, month, day, project", mode))
					}
				}
			}

			// project filter
			var projectIDs []string
			// resolve names or IDs to IDs only
			for _, nameOrID := range projectFilter {
				id := ""
				// if it's a name resolve it to the ID
				if project, err := ctx.Query.ProjectByFullName(strings.Split(nameOrID, "/")); err == nil {
					id = project.ID
				} else if _, err := ctx.Query.ProjectByID(nameOrID); err != nil {
					log.Fatal(fmt.Errorf("project %s not found", projectFilter))
				}
				projectIDs = append(projectIDs, id)
			}

			frameReport := report.NewBucketReport(model.NewSortedFrameList(ctx.Store.Frames()), ctx)
			frameReport.IncludeActiveFrames = includeActiveFrames
			frameReport.ProjectIDs = projectIDs
			frameReport.IncludeSubprojects = true
			frameReport.FilterRange = filterRange
			frameReport.RoundFramesTo = roundFrames
			frameReport.RoundingModeFrames = frameRoundingMode
			frameReport.RoundingModeTotals = totalsRoundingNode
			frameReport.RoundTotalsTo = roundTotals
			frameReport.SplitOperations = splitOperations
			frameReport.ShowEmptyBuckets = true
			frameReport.Update()

			if jsonOutput {
				data, err := json.MarshalIndent(frameReport.Result, "", "  ")
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(data))
			} else {
				options := htmlreport.Options{
					DecimalDuration:  decimalDurations,
					TemplateName:     templateName,
					TemplateFilePath: templateFilePath,
				}
				if err := printTemplate(ctx, frameReport, options); err != nil {
					log.Fatal(fmt.Errorf("error while rendering: %s", err.Error()))
				}
			}
		},
	}

	cmd.Flags().StringVarP(&templateName, "template", "", "default", "Template to use for rendering. This may either be a full path to a template file or the name (without extension) of a template shipped with gotime.")

	templateAnnotations := make(map[string][]string)
	templateAnnotations[cobra.BashCompFilenameExt] = []string{"gohtml"}
	cmd.Flags().StringVarP(&templateFilePath, "template-file", "", "", "Custom gohtml template file to use for rendering. See the website for more details.")
	cmd.Flag("template-file").Annotations = templateAnnotations

	// fixme add defaults?
	// cmd.Flags().BoolVarP(&includeActiveFrames, "current", "c", false, "(Don't) Include currently running frame in report.")
	cmd.Flags().StringVarP(&fromDateString, "from", "f", "", "The date when the report should start.")
	cmd.Flags().StringVarP(&toDateString, "to", "t", "", "Optional end date")

	// cmd.Flags().BoolVarP(&showAll, "all", "", false, "Reports all activities.")
	cmd.Flags().IntVarP(&year, "year", "y", 0, "Filter on a specific year. 0 is the current year, -1 is last year, etc.")
	// cmd.Flag("year").NoOptDefVal = "0"
	cmd.Flags().IntVarP(&month, "month", "m", 0, "Filter on a given month. For example, 0 is the current month, -1 is last month, etc.")
	cmd.Flag("month").NoOptDefVal = "0"
	cmd.Flags().IntVarP(&day, "day", "d", 0, "Select the date range of a given day. For example, 0 is today, -1 is one day ago, etc.")
	// cmd.Flag("day").NoOptDefVal = "0"

	cmd.Flags().StringSliceVarP(&projectFilter, "project", "p", []string{}, "--project ID | NAME . Reports activities only for the given project. You can add other projects by using this option multiple times.")

	cmd.Flags().StringVarP(&splitModes, "split", "s", "project", "Split the report into groups. Multiple values are possible. Possible values: year,month,week,day,project")

	cmd.Flags().DurationVarP(&roundFrames, "round-frames-to", "", time.Minute, "Round durations of each frame to the nearest multiple of this duration")
	cmd.Flags().StringVarP(&roundModeFrames, "round-frames", "", "", "Rounding mode for sums of durations. Default: no rounding. Possible values: up|nearest")

	cmd.Flags().DurationVarP(&roundTotals, "round-totals-to", "", time.Minute, "Round durations of each frame to the nearest multiple of this duration")
	cmd.Flags().StringVarP(&roundModeTotal, "round-totals", "", "", "Rounding mode for sums of durations. Default: no rounding. Possible values: up|nearest.")

	cmd.Flags().BoolVarP(&decimalDurations, "decimal", "", false, "Print durations as decimals 1.5h instead of 1:30h")

	cmd.Flags().BoolVarP(&jsonOutput, "json", "", false, "Prints JSON instead of plain text")

	// fixme
	// cmd.Flags().DurationVarP(&roundTotals, "round-totals-to", "", time.Duration(0), "Round the overall duration of each project to the next matching multiple of this duration")
	// cmd.Flags().StringVarP(&roundModeTotal, "round-totals", "", "up", "Rounding mode for sums of durations. Default: up. Possible values: up|nearest")

	parent.AddCommand(cmd)
	return cmd
}

func printTemplate(ctx *context.TomContext, report *report.BucketReport, opts htmlreport.Options) error {
	dir, _ := os.Getwd()

	t := htmlreport.NewReport(dir, opts, ctx)
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

func printReport(report *report.ResultBucket, ctx *context.TomContext, level int) {
	title := report.Title()
	if title != "" {
		printlnIndenting(level-1, title)
	} else if level == 1 {
		printlnIndenting(level-1, "Overall")
	}

	if !report.DateRange().Empty() {
		printfIndenting(level, "Date range: %s\n", report.DateRange().MinimalString())
	}
	if !report.TrackedDateRange().Empty() {
		printfIndenting(level, "Tracked dates: %s\n", report.TrackedDateRange().ShortString())
	}
	printfIndenting(level, "Duration: %s\n", ctx.DurationPrinter.Short(report.Duration.Get()))
	printfIndenting(level, "Exact Duration: %s\n", ctx.DurationPrinter.Short(report.Duration.GetExact()))
	fmt.Println()

	for _, r := range report.ChildBuckets {
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
