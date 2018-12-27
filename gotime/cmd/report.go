package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"../context"
	"../report"
)

func newReportCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var fromDate string
	var toDate string
	var roundFrames time.Duration
	var roundNearest bool
	var roundTotal time.Duration
	var roundTotalNearest bool

	var cmd = &cobra.Command{
		Use:   "report",
		Short: "Reporting about the tracked time",
		Run: func(cmd *cobra.Command, args []string) {
			frames := context.Store.Frames()
			if fromDate != "" {
				from, err := time.Parse(time.RFC3339Nano, fromDate)
				if err != nil {
					fatal(err)
				}

				for i, frame := range frames {
					if frame.Start != nil && frame.Start.Before(from) || frame.End != nil && frame.End.After(from) {
						frames = append(frames[:i], frames[i+1:]...)
					}
				}
			}

			if toDate != "" {
				to, err := time.Parse(time.RFC3339Nano, toDate)
				if err != nil {
					fatal(err)
				}

				for i, frame := range frames {
					if frame.Start != nil && frame.Start.Before(to) || frame.End != nil && frame.End.After(to) {
						frames = append(frames[:i], frames[i+1:]...)
					}
				}
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

			result, err := frameReport.Calc(frames, context)
			if err != nil {
				fatal(err)
			}

			if context.JsonOutput {
				data, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					fatal(err)
				}
				fmt.Println(string(data))
			} else {
				for _, r := range result {
					fmt.Printf("%s: %s\n", r.Name, r.Duration.String())
				}
			}
		},
	}

	cmd.Flags().StringVarP(&fromDate, "from", "f", "", "Optional start date")
	cmd.Flags().StringVarP(&toDate, "to", "t", "", "Optional end date")

	cmd.Flags().DurationVarP(&roundFrames, "round-frames", "r", time.Duration(0), "Round durations of each frame to the nearest multiple of this duration")
	cmd.Flags().BoolVarP(&roundNearest, "round-frames-nearest", "", false, "Round the durations of each frame to the nearest multiple. The default is to round up.")

	cmd.Flags().DurationVarP(&roundTotal, "round-total", "", time.Duration(0), "Round the overall duration of each project to the next matching multiple of this duration")
	cmd.Flags().BoolVarP(&roundTotalNearest, "round-total-nearest", "", false, "Round the overall duration of a project or tag the nearest multiple. The default is to round up.")

	parent.AddCommand(cmd)
	return cmd
}
