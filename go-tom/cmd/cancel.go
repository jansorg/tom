package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/cmd/util"
	"github.com/jansorg/tom/go-tom/context"
)

func newCancelCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	cancelAll := false

	var cmd = &cobra.Command{
		Use:   "cancel",
		Short: "cancel removed the currently running activity. No data will be recorded.",
		Run: func(cmd *cobra.Command, args []string) {
			frames := ctx.Query.ActiveFrames()
			if cancelAll {
				for _, f := range frames {
					if err := ctx.Store.RemoveFrame(f.ID); err != nil {
						util.Fatal(err)
					}
				}
				fmt.Printf("Successfully removed %d frames", len(frames))
			} else if len(frames) > 0 {
				sort.SliceStable(frames, func(i, j int) bool {
					return frames[i].IsBefore(frames[j])
				})
				if err := ctx.Store.RemoveFrame(frames[0].ID); err != nil {
					util.Fatal(err)
				}
				fmt.Println("Successfully stopped frame")
			} else {
				fmt.Println("no active frame found")
			}
		},
	}

	parent.AddCommand(cmd)
	return cmd
}
