package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/store"
)

func newStartCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	createOnTheFly := false
	allowMultiple := false
	var notes string

	var cmd = &cobra.Command{
		Use:   "start <project> [+tag1 +tag2]",
		Short: "starts a new timer",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			// tags := args[1:]

			projects := context.Query.ProjectsByShortNameOrID(projectName)
			if len(projects) > 1 {
				fatal(fmt.Errorf("more than one project found for %s", projectName))
			}

			var project *store.Project
			if len(projects) == 1 {
				project = projects[0]
			} else {
				if createOnTheFly == false {
					fatal(fmt.Errorf("project %s not found, on-the-fly is disabled", projectName))
				}
				var err error
				if project, err = context.Store.AddProject(store.Project{Name: projectName}); err != nil {
					fatal(err)
				}
			}

			frame := store.NewStartedFrame(project)
			frame.Notes = notes
			if _, err := context.Store.AddFrame(frame); err != nil {
				fatal(err)
			}
		},
	}

	cmd.Flags().BoolVarP(&createOnTheFly, "on-the-fly", "", true, "Create unknown projects and tags on the fly if set to true")
	viper.BindPFlag("on-the-fly", cmd.Flags().Lookup("on-the-fly"))
	viper.SetDefault("on-the-fly", false)

	cmd.Flags().BoolVarP(&allowMultiple, "allow-multiple", "", false, "Allow multiple active timers at the same time")
	viper.BindPFlag("allow-multiple", cmd.Flags().Lookup("allow-multiple"))
	viper.SetDefault("allow-multiple", "false")

	cmd.Flags().StringVarP(&notes, "notes", "", "", "Optional notes for the new time frame")

	parent.AddCommand(cmd)
	return cmd
}
