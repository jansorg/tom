package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"../context"
)

func newFramesCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "frames",
		Short: "Print a listing of all frames",
		Run: func(cmd *cobra.Command, args []string) {
			frames := context.Store.Frames()

			if context.JsonOutput {
				if bytes, err := json.MarshalIndent(frames, "", "  "); err != nil {
					fatal(err)
				} else {
					fmt.Println(string(bytes))
				}
			} else {
				for _, p := range frames {
					fmt.Println(p.Id)
				}
			}
		},
	}

	parent.AddCommand(cmd)
	return cmd
}
