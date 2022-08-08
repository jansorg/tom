package imports

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dataImport/tomImport"
	"github.com/jansorg/tom/go-tom/util"
)

func newTomImportCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "tom",
		Short: "import projects, tags and frames from a Tom data directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			absPath, err := filepath.Abs(args[0])
			if err != nil {
				util.Fatal(err)
			}

			if result, err := tomImport.NewImporter().Import(absPath, ctx); err != nil {
				util.Fatal(err)
			} else {
				fmt.Println(result.String())
			}
		},
	}

	parent.AddCommand(cmd)
	return cmd
}
