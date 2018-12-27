package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"../context"
	"../store"
)

func newStopCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	all := false
	notes := ""

	var cmd = &cobra.Command{
		Use:   "stop",
		Short: "stops the newest active timer. If --all is specified, then all active timers are stopped.",
		Run: func(cmd *cobra.Command, args []string) {
			active := context.Store.FindFrames(func(f store.Frame) bool {
				return f.IsActive()
			})

			sort.SliceStable(active, func(i, j int) bool {
				return active[i].Start.After(*active[j].Start)
			})

			if !all && len(active) > 0 {
				active = active[:1]
			}

			for _, frame := range active {
				frame.Stop()
				if notes != "" {
					frame.Notes = notes
				}
				if _, err := context.Store.UpdateFrame(frame); err != nil {
					fatal(err)
				}
			}

			fmt.Printf("Stopped %d timers\n", len(active))
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Stops all running timers, not just the newest")
	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Optional notes to set for all stopped timers")

	parent.AddCommand(cmd)
	return cmd
}
