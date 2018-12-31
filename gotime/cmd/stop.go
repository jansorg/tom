package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/activity"
	"github.com/jansorg/gotime/gotime/context"
)

func newStopCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	all := false
	notes := ""
	var tags []string

	var cmd = &cobra.Command{
		Use:   "stop",
		Short: "stops the newest active timer. If --all is specified, then all active timers are stopped.",
		Run: func(cmd *cobra.Command, args []string) {
			a := activity.NewActivityControl(context, false, false)

			count := 0
			if all {
				if frames, err := a.StopAll(notes, []string{}); err != nil {
					fatal(err)
				} else {
					count = len(frames)
				}
			} else {
				_, err := a.StopNewest(notes, tags)
				if err != nil {
					fatal(err)
				} else {
					count = 1
				}
			}

			fmt.Printf("Stopped %d timers\n", count)
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Stops all running activities, not just the newest")
	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Optional notes to set for all stopped activities")
	cmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "Optional tags to add to all stopped activities")

	parent.AddCommand(cmd)
	return cmd
}
