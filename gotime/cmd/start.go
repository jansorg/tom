package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jansorg/gotime/gotime/activity"
	"github.com/jansorg/gotime/gotime/context"
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
			a := activity.NewActivityControl(context, createOnTheFly, allowMultiple)

			_, err := a.Start(args[0], "", args[1:])
			if err != nil {
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
