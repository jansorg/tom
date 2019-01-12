package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/invoice/sevdesk"
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
				Fatal(err)
			}

			invoiceData, err := cfg.createSummary()
			if err != nil {
				Fatal(err)
			}

			if cfg.dryRun {
				fmt.Printf("Date range: %s\n", cfg.filterRange.MinimalString())
				for _, line := range invoiceData.lines {
					fmt.Printf("%s: %.2f hours at %.2f %s\n", line.ProjectName, line.Hours, line.HourlyRate, line.Currency)
				}
			} else {
				client := sevdesk.NewClient(apiKey)
				err = client.LoadBasicData()
				if err != nil {
					Fatal(err)
				}

				// try to find the contact
				contacts, err := client.GetContacts()
				if err != nil {
					Fatal(err)
				}


				var contactID string
				for _, contact := range contacts {
					found := strings.Contains(contact.Description, fmt.Sprintf("[gotime: %s]", invoiceData.projectID ))
					if found {
						contactID = contact.ID
						break
					}
				}

				if contactID == ""{
					// create a new company contact where invoices to this project will attach
					contact, err := client.CreateCompanyContact(client.NewCompanyContact(fmt.Sprintf("[gotime] Project: %s", invoiceData.projectName), fmt.Sprintf("[gotime: %s]", invoiceData.projectID)))
					if err != nil {
						Fatal(err)
					}
					contactID = contact.ID
				}

				// fixme
				invoice, err := client.NewInvoice(
					time.Now(),
					fmt.Sprintf("%s %s", cfg.project.FullName, cfg.filterRange.MinimalString()),
					contactID,
					100,
					invoiceData.taxRate,
					"",
					sevdesk.TaxTypeNotEU,
					sevdesk.Currency(invoiceData.currency),
					0,
					invoiceData.address)

				if err != nil {
					Fatal(err)
				}

				resp, err := client.CreateInvoice(invoice)
				if err != nil {
					Fatal(err)
				}

				for _, line := range invoiceData.lines {
					posDef, err := client.NewInvoicePosition(resp.ID, line.ProjectName, line.Hours, "hours", line.HourlyRate, 0)
					if err != nil {
						Fatal(err)
					}

					_, err = client.CreateInvoicePos(posDef)
					if err != nil {
						Fatal(err)
					}
				}

				fmt.Printf("Successfully created invoice with %d positions. URL: %s\n", len(invoiceData.lines), resp.BrowserURL())
			}
		},
	}

	cmd.Flags().StringVarP(&apiKey, "key", "k", "", "The API key to use for sevdesk.com")
	cmd.MarkFlagRequired("key")

	parent.AddCommand(cmd)
	return cmd
}
