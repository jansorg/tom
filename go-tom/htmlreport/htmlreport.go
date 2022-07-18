package htmlreport

import (
	"bytes"
	"fmt"
	htmlTemplate "html/template"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"time"

	assetTemplate "github.com/arschles/go-bindata-html-template"
	"golang.org/x/text/message"

	"github.com/jansorg/tom/go-tom"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateTime"
	"github.com/jansorg/tom/go-tom/money"
	"github.com/jansorg/tom/go-tom/report"
)

type Report struct {
	workingDir string
	options    Options
	ctx        *context.TomContext
}

type Options struct {
	CustomTitle        *string          `json:"title"`
	CustomDescription  *string          `json:"description"`
	ShowSummary        bool             `json:"show_summary"`
	ShowExactDurations bool             `json:"show_exact"`
	ShowMatrixTables   bool             `json:"matrix_tables"`
	DecimalDuration    bool             `json:"decimal_duration"`
	ShowSales          bool             `json:"show_sales"`
	ShowTracked        bool             `json:"show_tracked"`
	ShowUnTracked      bool             `json:"show_untracked"`
	TemplateName       *string          `json:"template_name"`
	TemplateFilePath   *string          `json:"template_path"`
	CustomCSS          htmlTemplate.CSS `json:"css"`
	CustomCSSFile      string           `json:"css_file"`

	Report report.Config `json:"report"`
}

var DefaultOptions = Options{
	ShowMatrixTables: true,
}

func NewReport(workingDir string, opts Options, ctx *context.TomContext) *Report {
	return &Report{
		options:    opts,
		workingDir: workingDir,
		ctx:        ctx,
	}
}

func (r *Report) Render(results *report.BucketReport) ([]byte, error) {
	functionMap := map[string]interface{}{
		"reportOptions": func() *Options {
			return &r.options
		},
		"add": func(a, b int) int {
			return a + b
		},
		"i18n": func(key string) string {
			return r.ctx.LocalePrinter.Sprintf(key)
		},
		"inlineCSS": func(filename string) htmlTemplate.CSS {
			file, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Println(fmt.Errorf("error reading CSS file %s", filename))
				return htmlTemplate.CSS(fmt.Sprintf("/* file not found: %s*/", filename))
			}
			return htmlTemplate.CSS(file)
		},
		"langBase": func() string {
			base, _ := r.ctx.Language.Base()
			return base.String()
		},
		"formatNumber": func(n interface{}) string {
			if floatValue, ok := n.(float64); ok {
				return r.ctx.LocalePrinter.Sprintf(message.Key("float-format", "%.2f"), floatValue)
			}
			if floatValue, ok := n.(float32); ok {
				return r.ctx.LocalePrinter.Sprintf(message.Key("float-format", "%.2f"), floatValue)
			}
			return r.ctx.LocalePrinter.Sprint(n)
		},
		"formatTime": func(date time.Time) string {
			if r.showSeconds() {
				return r.ctx.Locale.FmtTimeMedium(date)
			}
			return r.ctx.Locale.FmtTimeShort(date)
		},
		"formatDate": func(date time.Time) string {
			return r.ctx.Locale.FmtDateShort(date)
		},
		"formatDateTime": func(date time.Time) string {
			return r.ctx.DateTimePrinter.DateTime(date)
		},
		"formatMoney": func(money *money.Money) string {
			if money == nil {
				return ""
			}
			// fixme make i18n aware
			return money.Currency().Formatter().Format(money.Amount())
		},
		"roundedDuration": func(duration time.Duration, bucket report.ResultBucket) time.Duration {
			return bucket.Duration.CalculateRoundedDuration(duration)
		},
		"minDuration": func(duration time.Duration) string {
			if r.options.DecimalDuration {
				return r.ctx.DecimalDurationPrinter.Minimal(duration, r.showSeconds())
			}
			return r.ctx.DurationPrinter.Minimal(duration, r.showSeconds())
		},
		"shortDuration": func(duration time.Duration) string {
			if r.options.DecimalDuration {
				return r.ctx.DecimalDurationPrinter.Short(duration, r.showSeconds())
			}
			return r.ctx.DurationPrinter.Short(duration, r.showSeconds())
		},
		"longDuration": func(duration time.Duration) string {
			if r.options.DecimalDuration {
				return r.ctx.DecimalDurationPrinter.Long(duration, r.showSeconds())
			}
			return r.ctx.DurationPrinter.Long(duration, r.showSeconds())
		},
		"isMatrix": report.IsMatrix,
		"sumChildValues": func(parent *report.ResultBucket, childIndex int) *dateTime.DurationSum {
			sum := dateTime.NewDurationSum()
			if parent == nil {
				return sum
			}

			for _, b := range parent.ChildBuckets {
				childCount := len(b.ChildBuckets)
				if childCount >= childIndex && childIndex < childCount {
					sum.AddSum(b.ChildBuckets[childIndex].Duration)
				}
			}

			return sum
		},
		"safeHTML": func(html string) htmlTemplate.HTML {
			return htmlTemplate.HTML(html)
		},
	}

	if r.options.TemplateFilePath != nil && *r.options.TemplateFilePath != "" {
		templatePath := *r.options.TemplateFilePath
		if !filepath.IsAbs(templatePath) {
			templatePath = filepath.Join(r.workingDir, templatePath)
		}

		files, err := filepath.Glob(filepath.Join(filepath.Dir(templatePath), "*.gohtml"))

		tmpl, err := htmlTemplate.New(filepath.Base(templatePath)).Funcs(functionMap).ParseFiles(append(files, templatePath)...)
		if err != nil {
			return nil, err
		}
		out := bytes.NewBuffer(nil)
		if err = tmpl.Execute(out, results); err != nil {
			return nil, err
		}

		return out.Bytes(), nil
	} else if r.options.TemplateName != nil && *r.options.TemplateName != "" {
		templatePath := *r.options.TemplateName

		baseDir := path.Join("reports", "html")
		templateFiles := []string{
			path.Join(baseDir, templatePath+".gohtml"),
			path.Join(baseDir, "commons.gohtml"),
		}

		tmpl, err := assetTemplate.New(templatePath, tom.Asset).Funcs(functionMap).ParseFiles(templateFiles...)
		if err != nil {
			return nil, err
		}
		out := bytes.NewBuffer(nil)
		if err = tmpl.Execute(out, results); err != nil {
			return nil, err
		}

		return out.Bytes(), nil
	} else {
		return nil, fmt.Errorf("template undefined")
	}
}

// showSeconds returns if durations and timestamps should display seconds in the report
func (r *Report) showSeconds() bool {
	return r.options.Report.EntryRounding.IsSecondPrecision()
}
