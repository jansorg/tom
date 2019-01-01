package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/invoice/sevdesk"
)

func newSevdeskCommand(ctx *context.GoTimeContext, parent *cobra.Command) *cobra.Command {
	apiKey := ""

	var cmd = &cobra.Command{
		Use:     "sevdesk",
		Short:   "Create a new invoice at sevdesk.com",
		Example: "gotime invoice sevdesk --project myProject",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := parseInvoiceCmd(ctx, cmd)
			if err != nil {
				fatal(err)
			}

			lines, err := cfg.createSummary()
			if err != nil {
				fatal(err)
			}

			if cfg.dryRun {
				fmt.Printf("Date range: %s\n", cfg.filterRange.MinimalString())
				for _, line := range lines {
					fmt.Printf("%.2f at %.2f %s %s", line.Hours, line.HourlyRate, line.Currency, line.ProjectName)
				}
			} else {
				client := sevdesk.NewClient(apiKey)
				err = client.LoadBasicData()
				if err != nil {
					fatal(err)
				}

				invoice, err := client.NewInvoice(time.Now(), fmt.Sprintf("%s %s", cfg.project.FullName, cfg.filterRange.MinimalString()), "7067576", 100, 0.0, "", sevdesk.TaxTypeNotEU, sevdesk.USD, 0, "")
				if err != nil {
					fatal(err)
				}

				resp, err := client.CreateInvoice(invoice)
				if err != nil {
					fatal(err)
				}

				for _, line := range lines {
					posDef, err := client.NewInvoicePosition(resp.ID, line.ProjectName, line.Hours, "hours", line.HourlyRate, 0)
					if err != nil {
						fatal(err)
					}

					_, err = client.CreateInvoicePos(posDef)
					if err != nil {
						fatal(err)
					}
				}

				fmt.Printf("Successfully created invoice with %d positions. URL: %s\n", len(lines), resp.BrowserURL())
			}
		},
	}

	cmd.Flags().StringVarP(&apiKey, "key", "k", "", "The API key to use for sevdesk.com")
	cmd.MarkPersistentFlagRequired("key")

	parent.AddCommand(cmd)
	return cmd
}
