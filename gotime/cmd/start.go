package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jansorg/gotime/gotime/activity"
	"github.com/jansorg/gotime/gotime/config"
	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/store"
)

func newStartCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var notes string
	var shiftedStart time.Duration

	var cmd = &cobra.Command{
		Use:   "start [--past <duration>] <project> [+tag1 +tag2]",
		Short: "starts a new activity for the given project ands adds a list of optional tags",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createMissingProject := viper.GetBool(config.KeyProjectCreateMissing)
			stopActives := viper.GetBool(config.KeyActivityStopOnStart)

			control := activity.NewActivityControl(ctx, createMissingProject, false, time.Now().Add(shiftedStart))

			tags, err := argsToTags(ctx, args[1:])
			if err != nil {
				fatal(err)
			}

			var stoppedFrames []*store.Frame
			if stopActives {
				// fixme tags for stop?
				stoppedFrames, _ = control.StopAll("", nil)
			}

			frame, err := control.Start(args[0], "", tags)
			if err == activity.ProjectNotFoundErr {
				fatal(fmt.Errorf("project %s not found. Use --create-missing to create missing projects on-the-fly", args[0]))
			} else if err != nil {
				fatal(err)
			}

			if stoppedFrames != nil {
				// fixme i18n?
				fmt.Printf("Stopped %d activities\n", len(stoppedFrames))
			}

			if project, err := ctx.Query.ProjectByID(frame.ProjectId); err == nil {
				// fixme i18n?
				fmt.Printf("Started new activity for %s at %v. Tags: %s\n", project.FullName, ctx.DateTimePrinter.Time(*frame.Start), args[1:])
			}
		},
	}

	cmd.Flags().Bool("create-missing", false, "")
	if err := viper.BindPFlag(config.KeyProjectCreateMissing, cmd.Flag("create-missing")); err != nil {
		fatal(err)
	}

	cmd.Flags().Bool("stop-on-start", false, "")
	if err := viper.BindPFlag(config.KeyActivityStopOnStart, cmd.Flag("stop-on-start")); err != nil {
		fatal(err)
	}

	cmd.Flags().DurationVarP(&shiftedStart, "past", "d", 0, "Duration to add to the new activity. The current activity will be started in the past. Running activities will be stopped before the new activity.")
	cmd.Flags().StringVarP(&notes, "notes", "", "", "Optional notes for the new time frame")

	parent.AddCommand(cmd)
	return cmd
}
