package remove

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/util"
)

func newRemoveFrameCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "frame ID ...",
		Short: "removes one or more frames, identified by ID",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			removed := 0
			notFound := 0

			ctx.Store.StartBatch()
			for _, id := range args {
				if err := ctx.Store.RemoveFrame(id); err != nil {
					notFound++
				} else {
					removed++
				}
			}
			if err := ctx.Store.StopBatch(); err != nil {
				util.Fatal(err)
			}

			fmt.Printf("%d frames removed, %d frames not found.\n", removed, notFound)
		},
	}
	parent.AddCommand(cmd)

	return cmd
}
