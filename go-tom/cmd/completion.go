package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/go-tom/context"
)

func newCompletionCommand(context *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	// completionCmd represents the completion command
	var completionCmd = &cobra.Command{
		Use:    "completion",
		Hidden: true,
		Short:  "Generates bash completion scripts",
		Long: `To load completion run

. <(gotime completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(gotime completion)
`,
		Run: func(cmd *cobra.Command, args []string) {
			RootCmd.GenBashCompletion(os.Stdout);
		},
	}

	parent.AddCommand(completionCmd)
	return completionCmd
}
