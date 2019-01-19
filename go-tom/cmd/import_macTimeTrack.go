package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/import/macTimeTracker"
)

func newImportMacTimeTrackCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "macTimeTracker timeTrackerExport.csv",
		Short: "",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			created, err := macTimeTracker.ImportCSV(args[0], ctx)
			if err != nil {
				Fatal(err)
			} else {
				fmt.Printf("created %d frames\n", created)
			}
		},
	}

	parent.AddCommand(cmd)
	return cmd
}
