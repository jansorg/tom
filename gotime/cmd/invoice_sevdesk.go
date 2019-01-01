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

			ivConfig, err := cfg.createSummary()
			if err != nil {
				fatal(err)
			}

			if cfg.dryRun {
				fmt.Printf("Date range: %s\n", cfg.filterRange.MinimalString())
				for _, line := range ivConfig.lines {
					fmt.Printf("%s: %.2f hours at %.2f %s\n", line.ProjectName, line.Hours, line.HourlyRate, line.Currency)
				}
			} else {
				client := sevdesk.NewClient(apiKey)
				err = client.LoadBasicData()
				if err != nil {
					fatal(err)
				}

				// fixme
				invoice, err := client.NewInvoice(time.Now(),
					fmt.Sprintf("%s %s", cfg.project.FullName, cfg.filterRange.MinimalString()),
					"7067576",
					100,
					ivConfig.taxRate,
					"",
					sevdesk.TaxTypeNotEU,
					sevdesk.Currency(ivConfig.currency),
					0,
					ivConfig.address)

				if err != nil {
					fatal(err)
				}

				resp, err := client.CreateInvoice(invoice)
				if err != nil {
					fatal(err)
				}

				for _, line := range ivConfig.lines {
					posDef, err := client.NewInvoicePosition(resp.ID, line.ProjectName, line.Hours, "hours", line.HourlyRate, 0)
					if err != nil {
						fatal(err)
					}

					_, err = client.CreateInvoicePos(posDef)
					if err != nil {
						fatal(err)
					}
				}

				fmt.Printf("Successfully created invoice with %d positions. URL: %s\n", len(ivConfig.lines), resp.BrowserURL())
			}
		},
	}

	cmd.Flags().StringVarP(&apiKey, "key", "k", "", "The API key to use for sevdesk.com")
	cmd.MarkFlagRequired("key")

	parent.AddCommand(cmd)
	return cmd
}
