package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/go-tom/activity"
	"github.com/jansorg/gotime/go-tom/context"
	"github.com/jansorg/gotime/go-tom/store"
)

func newStopCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	all := false
	notes := ""
	// var tags []string
	var shiftedStop time.Duration

	var cmd = &cobra.Command{
		Use:   "stop [--past <duration>] [--all] [--notes \"override notes\"]",
		Short: "stops the newest active timer. If --all is specified, then all active timers are stopped.",
		Run: func(cmd *cobra.Command, args []string) {
			a := activity.NewActivityControl(ctx, false, false, time.Now().Add(shiftedStop))

			tags, err := argsToTags(ctx, args)
			if err != nil {
				fatal(err)
			}

			var frames []*store.Frame
			if all {
				if frames, err = a.StopAll(notes, tags); err != nil {
					fatal(err)
				}
			} else {
				frame, err := a.StopNewest(notes, tags)
				if err != nil {
					fatal(err)
				}
				frames = []*store.Frame{frame}
			}

			// translate
			fmt.Printf("Stopped %d timers at %s\n", len(frames), ctx.DateTimePrinter.Time(time.Now()))
			for _, frame := range frames {
				fmt.Printf("\t%s\n", ctx.DurationPrinter.Minimal(frame.Duration()))
			}
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Stops all running activities, not just the newest")
	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Optional notes to set for all stopped activities")
	// cmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "Optional tags to add to all stopped activities")
	cmd.Flags().DurationVarP(&shiftedStop, "past", "d", 0, "Stop the activity this duration before now, e.g. `--past 5m` stops the activity 5m before the current time")

	parent.AddCommand(cmd)
	return cmd
}
