package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jansorg/gotime/gotime/activity"
	"github.com/jansorg/gotime/gotime/config"
	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/store"
)

func newStartCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var notes string

	var cmd = &cobra.Command{
		Use:   "start <project> [+tag1 +tag2]",
		Short: "starts a new activity for the given project ands adds a list of optional tags",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createMissingProject := viper.GetBool(config.KeyProjectCreateMissing)
			stopActives := viper.GetBool(config.KeyActivityStopOnStart)

			control := activity.NewActivityControl(context, createMissingProject, false)

			var stoppedFrames []*store.Frame
			if stopActives {
				stoppedFrames, _ = control.StopAll("", []string{})
			}

			frame, err := control.Start(args[0], "", args[1:])
			if err == activity.ProjectNotFoundErr {
				fatal(fmt.Errorf("Project not found. Use --create-missing to create missing projects on-the-fly."))
			} else if err != nil {
				fatal(err)
			}

			if stoppedFrames != nil {
				// fixme i18n?
				fmt.Printf("Stopped %d activities\n", len(stoppedFrames))
			}

			if project, err := context.Query.ProjectByID(frame.ProjectId); err == nil {
				// fixme i18n?
				fmt.Printf("Started new activity for %s at %v\n", project.FullName, ctx.DateTimePrinter.Time(*frame.Start))
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

	cmd.Flags().StringVarP(&notes, "notes", "", "", "Optional notes for the new time frame")

	parent.AddCommand(cmd)
	return cmd
}
