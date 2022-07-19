package edit

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/util"
)

func newEditFrameCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var startTime string
	var endTime string
	var projectIDOrName string
	var nameDelimiter string
	var notes string
	var archive bool

	var cmd = &cobra.Command{
		Use:   "frame ID",
		Short: "edit properties of a frame",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var usedStart, usedEnd, usedProjectID, usedNotes *string
			var archiveFrames *bool

			if !cmd.Flag("start").Changed {
				usedStart = nil
			} else {
				usedStart = &startTime
			}

			if !cmd.Flag("end").Changed {
				usedEnd = nil
			} else {
				usedEnd = &endTime
			}

			if !cmd.Flag("notes").Changed {
				usedNotes = nil
			} else {
				usedNotes = &notes
			}

			if !cmd.Flag("project").Changed {
				usedProjectID = nil
			} else {
				usedProjectID = &projectIDOrName
			}

			if !cmd.Flag("archived").Changed {
				archiveFrames = nil
			} else {
				archiveFrames = &archive
			}

			if err := doEditFrameCommand(ctx, args, usedStart, usedEnd, usedNotes, usedProjectID, nameDelimiter, archiveFrames); err != nil {
				util.Fatal(err)
			} else {
				fmt.Println("successfully updated")
			}
		},
	}

	cmd.Flags().StringVarP(&startTime, "start", "f", "", "update the start time.")
	cmd.Flags().StringVarP(&endTime, "end", "t", "", "update the end time. Pass an empty value to remove the end time.")
	cmd.Flags().StringVarP(&notes, "notes", "n", "", "updates the notes for the given frame. Pass an empty string to remove the notes from the frame.")
	cmd.Flags().StringVarP(&projectIDOrName, "project", "p", "", "Project ID or full name to use as new project for all passed frame IDs")
	cmd.Flags().StringVarP(&nameDelimiter, "name-delimiter", "", "/", "Delimiter used in full project names")
	cmd.Flags().BoolVarP(&archive, "archived", "", archive, "Sets the archived flag")

	parent.AddCommand(cmd)
	return cmd
}

func doEditFrameCommand(ctx *context.TomContext, frameIDs []string, startTime, endTime, notes, projectIDOrName *string, nameDelimiter string, archived *bool) error {
	// make sure that all frames exist before applying updates
	frames, err := ctx.Query.FramesByID(frameIDs...)
	if err != nil {
		return err
	}

	// validate project
	validatedProjectID := ""
	if projectIDOrName != nil && *projectIDOrName != "" {
		if p, err := ctx.Query.ProjectByFullNameOrID(*projectIDOrName, nameDelimiter); err != nil {
			return err
		} else {
			validatedProjectID = p.ID
		}
	}

	ctx.Store.StartBatch()
	defer ctx.Store.StopBatch()

	for _, frame := range frames {
		if startTime != nil {
			if *startTime == "" {
				return fmt.Errorf("empty start time is not allowed")
			} else if start, err := time.Parse(time.RFC3339, *startTime); err != nil {
				return err
			} else {
				frame.Start = &start
			}
		}

		if endTime != nil {
			if *endTime == "" {
				frame.End = nil
			} else if end, err := time.Parse(time.RFC3339, *endTime); err != nil {
				return err
			} else {
				frame.End = &end
			}
		}

		if projectIDOrName != nil {
			frame.ProjectId = validatedProjectID
		}

		if notes != nil {
			frame.Notes = *notes
		}

		if archived != nil {
			frame.Archived = *archived
		}

		if frame, err = ctx.Store.UpdateFrame(*frame); err != nil {
			return err
		}
	}

	return nil
}
