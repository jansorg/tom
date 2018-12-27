package cmd

import (
	"sort"

	"github.com/spf13/cobra"

	"../store"
)

func newStopCommand(context *GoTimeContext, parent *cobra.Command) *cobra.Command {
	all := false

	var cmd = &cobra.Command{
		Use:   "stop",
		Short: "stops the newest active timer. If --all is specified, then all active timers are stopped.",
		Run: func(cmd *cobra.Command, args []string) {
			active := context.Store.FindFrames(func(f store.Frame) bool {
				return f.IsActive()
			})

			if all {
				for _, frame := range active {
					frame.Stop()
					if _, err := context.Store.UpdateFrame(frame); err != nil {
						fatal(err)
					}
				}
			} else {
				sort.SliceStable(active, func(i, j int) bool {
					return active[i].Start.After(active[j].Start)
				})
				last := active[0]
				last.Stop()
				if _, err := context.Store.UpdateFrame(last); err != nil {
					fatal(err)
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Stops all running timers, not just the newest")

	parent.AddCommand(cmd)
	return cmd
}
