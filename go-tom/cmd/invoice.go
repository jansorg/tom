package cmd

import (
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/report"
	"github.com/jansorg/tom/go-tom/util"
)

func newInvoiceCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "invoice",
		Short: "Create a new invoice in a cloud based service. See the list of command to see which service APIs are supported at this time.",
	}

	cmd.PersistentFlags().BoolP("dry-run", "d", false, "Dry run without creating data on the remote side")
	cmd.PersistentFlags().StringP("project", "p", "", "Project to list in the invoice.")

	cmd.PersistentFlags().DurationP("round-frames-to", "", 0, "Round frames to multiples of this duration. Default: 0 (no rounding)")
	cmd.PersistentFlags().StringP("round-frames", "", "up", "Rounding mode for time frames of the projects. Default: up (always round up)")

	cmd.PersistentFlags().IntP("month", "", -1, "The month to list in the reporting. Default: -1 (last month)")
	// cmd.PersistentFlags().StringP("from", "", "", "Start time for frames to list in the invoice.")
	// cmd.PersistentFlags().StringP("to", "", "", "End time for frames to list in the invoice")

	cmd.MarkPersistentFlagRequired("project")

	newSevdeskCommand(ctx, cmd)
	parent.AddCommand(cmd)
	return cmd
}

type invoiceCmdConfig struct {
	ctx             *context.TomContext
	dryRun          bool
	project         *model.Project
	filterRange     util.DateRange
	roundFramesTo   time.Duration
	roundFramesMode util.RoundingMode
}

type invoiceConfig struct {
	currency    string
	taxRate     float64
	lines       []ProjectInvoiceLine
	address     string
	projectID   string
	projectName string
}

func (c invoiceCmdConfig) createSummary() (invoiceConfig, error) {
	storeFrames := c.ctx.Store.Frames()

	reportConfig := report.Config{
		ProjectIDs:         []string{c.project.ID},
		IncludeSubprojects: true,
		DateFilterRange:    c.filterRange,
		EntryRounding: util.RoundingConfig{
			Mode: c.roundFramesMode,
			Size: c.roundFramesTo,
		},
	}

	frameReport := report.NewBucketReport(model.NewSortedFrameList(storeFrames), reportConfig, c.ctx)
	result := frameReport.Update()

	// desc, _ := c.ctx.Query.GetInheritedStringProp(c.project.ID, config.InvoiceDescriptionProperty)
	// address, _ := c.ctx.Query.GetInheritedStringProp(c.project.ID, config.InvoiceAddressProperty)
	// currency, _ := c.ctx.Query.GetInheritedStringProp(c.project.ID, config.InvoiceCurrencyProperty)
	// hourlyRate, _ := c.ctx.Query.GetInheritedFloatProp(c.project.ID, config.InvoiceHourlyRateProperty)
	// taxRate, _ := c.ctx.Query.GetInheritedFloatProp(c.project.ID, config.InvoiceTaxRateProperty)

	return invoiceConfig{
		projectID:   c.project.ID,
		projectName: c.project.Name,
		// currency:    currency,
		// taxRate:     taxRate,
		// address:     address,
		lines: []ProjectInvoiceLine{
			{
				ProjectName: c.project.Name,
				Hours:       result.Duration.Get().Hours(),
				// Description: desc,
				// Currency:    currency,
				// HourlyRate:  hourlyRate,
			},
		},
	}, nil
}

type ProjectInvoiceLine struct {
	ProjectName string
	Description string
	Hours       float64
	HourlyRate  float64
	Currency    string
}

func parseInvoiceCmd(ctx *context.TomContext, cmd *cobra.Command) (invoiceCmdConfig, error) {
	var filterRange util.DateRange

	// fixme add start and end date

	if filterRange.Empty() {
		if month, err := cmd.Flags().GetInt("month"); err != nil {
			return invoiceCmdConfig{}, err
		} else {
			filterRange = util.NewMonthRange(time.Now(), ctx.Locale, time.Local).Shift(0, month, 0)
		}
	}

	projectName, err := cmd.Flags().GetString("project")
	if err != nil {
		return invoiceCmdConfig{}, err
	}

	project, err := ctx.Query.ProjectByFullName(strings.Split(projectName, "/"))
	if err != nil {
		return invoiceCmdConfig{}, err
	}

	roundModeFrames, err := cmd.Flags().GetString("round-frames")
	if err != nil {
		return invoiceCmdConfig{}, err
	}
	frameRoundingMode := util.ParseRoundingMode(roundModeFrames)

	roundFramesTo, err := cmd.Flags().GetDuration("round-frames-to")
	if err != nil {
		return invoiceCmdConfig{}, err
	}

	dry, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		return invoiceCmdConfig{}, err
	}

	return invoiceCmdConfig{
		ctx:             ctx,
		dryRun:          dry,
		project:         project,
		filterRange:     filterRange,
		roundFramesMode: frameRoundingMode,
		roundFramesTo:   roundFramesTo,
	}, nil
}
