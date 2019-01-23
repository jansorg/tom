package edit

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
)

func newEditFrameCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var startTime string
	var endTime string
	var notes string

	var cmd = &cobra.Command{
		Use:   "frame ID",
		Short: "edit properties of a frame",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			frame, err := ctx.Query.FrameByID(args[0])
			if err != nil {
				log.Fatal(err)
			}

			fromFlag := cmd.Flag("start")
			if fromFlag.Changed {
				if startTime == "" {
					log.Fatal(fmt.Errorf("empty start time is not allowed"))
				} else if start, err := time.Parse(time.RFC3339, startTime); err != nil {
					log.Fatal(err)
				} else {
					frame.Start = &start
				}
			}

			toFlag := cmd.Flag("end")
			if toFlag.Changed {
				if endTime == "" {
					frame.End = nil
				} else if end, err := time.Parse(time.RFC3339, endTime); err != nil {
					log.Fatal(err)
				} else {
					frame.End = &end
				}
			}

			notesFlag := cmd.Flag("notes")
			if notesFlag.Changed {
				frame.Notes = notes
			}

			if frame, err = ctx.Store.UpdateFrame(*frame); err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("successfully updated frame %s\n", frame.ID)
			}
		},
	}

	cmd.Flags().StringVarP(&startTime, "start", "f", "", "update the start time.")
	cmd.Flags().StringVarP(&endTime, "end", "t", "", "update the end time. Pass an empty value to remove the end time.")
	cmd.Flags().StringVarP(&notes, "notes", "n", "", "updates the notes for the given frame. Pass an empty string to remove the notes from the frame.")

	parent.AddCommand(cmd)
	return cmd
}
