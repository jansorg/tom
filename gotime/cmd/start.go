package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"../store"
)

func newStartCommand(context *GoTimeContext, parent *cobra.Command) *cobra.Command {
	createOnTheFly := false
	allowMultiple := false

	var cmd = &cobra.Command{
		Use:   "start <project> [+tag1 +tag2]",
		Short: "starts a new timer",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			// tags := args[1:]

			projects := context.Store.FindProjects(func(p store.Project) bool {
				return p.ShortName == projectName || p.Id == projectName
			})
			if len(projects) > 1 {
				fatal(fmt.Errorf("more than one project found for %s", projectName))
			}

			var project store.Project
			if len(projects) == 1 {
				project = projects[0]
			} else {
				if createOnTheFly == false {
					fatal(fmt.Errorf("project %s not found, on-the-fly is disabled", projectName))
				}
				var err error
				if project, err = context.Store.AddProject(store.Project{ShortName: projectName, FullName: projectName}); err != nil {
					fatal(err)
				}
			}

			if _, err := context.Store.AddFrame(store.NewStartedFrame(project)); err != nil {
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

	parent.AddCommand(cmd)
	return cmd
}
