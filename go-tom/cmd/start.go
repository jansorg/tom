package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jansorg/tom/go-tom/activity"
	"github.com/jansorg/tom/go-tom/config"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
)

func newStartCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	var notes string

	var cmd = &cobra.Command{
		Use:     "start <project> [time shift into past] [+tag1 +tag2]",
		Short:   "starts a new activity for the given project ands adds a list of optional tags",
		Example: "start acme 15m +onsite",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createMissingProject := viper.GetBool(config.KeyProjectCreateMissing)
			stopActives := viper.GetBool(config.KeyActivityStopOnStart)

			projectName := args[0]
			var shiftedStart time.Duration

			// look out for a time shift on the command line
			if len(args) >= 2 && !strings.HasPrefix(args[1], "+") /*&& !cmd.Flag("past").Changed */{
				if shift, err := time.ParseDuration(args[1]); err == nil {
					// it's not making sense to start a task in the future. Also, - is parsed as a shorthand flag prefix and we don't want the user working around that all the time
					if shift.Seconds() > 0 {
						shift = -shift
					}
					shiftedStart = shift
					args = args[2:]
				}
			} else {
				args = args[1:]
			}

			control := activity.NewActivityControl(ctx, createMissingProject, false, time.Now().Add(shiftedStart))

			tags, err := argsToTags(ctx, args)
			if err != nil {
				fatal(err)
			}

			var stoppedFrames []*model.Frame
			if stopActives {
				// fixme tags for stop?
				stoppedFrames, _ = control.StopAll("", nil)
			}

			frame, err := control.Start(projectName, "", tags)
			if err == activity.ProjectNotFoundErr {
				fatal(fmt.Errorf("project %s not found. Use --create-missing to create missing projects on-the-fly", projectName))
			} else if err != nil {
				fatal(err)
			}

			if stoppedFrames != nil {
				// fixme i18n?
				fmt.Printf("Stopped %d activities\n", len(stoppedFrames))
			}

			if project, err := ctx.Query.ProjectByID(frame.ProjectId); err == nil {
				// fixme i18n?
				fmt.Printf("Started new activity for %s at %v. Tags: %s\n", project.GetFullName("/"), ctx.DateTimePrinter.Time(*frame.Start), args)
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
