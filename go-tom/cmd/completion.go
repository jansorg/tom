package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
)

func newCompletionCommand(context *context.TomContext, parent *cobra.Command) *cobra.Command {
	// completionCmd represents the completion command
	var completionCmd = &cobra.Command{
		Use:    "completion",
		Hidden: true,
		Short:  "Generates bash completion scripts",
		Long: `To load completion run

. <(tom completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(tom completion)
`,
		Run: func(cmd *cobra.Command, args []string) {
			RootCmd.GenBashCompletion(os.Stdout);
		},
	}

	parent.AddCommand(completionCmd)
	return completionCmd
}
