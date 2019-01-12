package htmlreport

import (
	"bytes"
	"path"
	"time"

	"github.com/arschles/go-bindata-html-template"

	"github.com/jansorg/tom/go-tom"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/report"
)

type Report struct {
	workingDir   string
	templatePath string
	options      Options
	ctx          *context.GoTimeContext
}

type Options struct {
	DecimalDurationn bool
}

func NewReport(workingDir string, templatePath string, opts Options, ctx *context.GoTimeContext) *Report {
	return &Report{
		options:      opts,
		workingDir:   workingDir,
		templatePath: templatePath,
		ctx:          ctx,
	}
}

func (r *Report) Render(results *report.BucketReport) (string, error) {
	functionMap := map[string]interface{}{
		"i18n": func(key string) string {
			return r.ctx.LocalePrinter.Sprintf(key)
		},
		"langBase": func() string {
			base, _ := r.ctx.Language.Base()
			return base.String()
		},
		"formatNumber": func(n interface{}) string {
			return r.ctx.LocalePrinter.Sprint(n)
		},
		"formatTime": func(date time.Time) string {
			return r.ctx.DateTimePrinter.Time(date)
		},
		"formatDate": func(date time.Time) string {
			return r.ctx.DateTimePrinter.Date(date)
		},
		"formatDateTime": func(date time.Time) string {
			return r.ctx.DateTimePrinter.DateTime(date)
		},
		"minDuration": func(duration time.Duration) string {
			if r.options.DecimalDurationn {
				return r.ctx.DecimalDurationPrinter.Minimal(duration)
			}
			return r.ctx.DurationPrinter.Minimal(duration)
		},
		"shortDuration": func(duration time.Duration) string {
			if r.options.DecimalDurationn {
				return r.ctx.DecimalDurationPrinter.Short(duration)
			}
			return r.ctx.DurationPrinter.Short(duration)
		},
		"longDuration": func(duration time.Duration) string {
			if r.options.DecimalDurationn {
				return r.ctx.DecimalDurationPrinter.Long(duration)
			}
			return r.ctx.DurationPrinter.Long(duration)
		},
	}

	baseDir := path.Join("reports", "html")
	templateFiles := []string{
		path.Join(baseDir, r.templatePath+".gohtml"),
		path.Join(baseDir, "commons.gohtml"),
	}

	tmpl, err := template.New(r.templatePath, tom.Asset).Funcs(functionMap).ParseFiles(templateFiles...)
	if err != nil {
		return "", err
	}

	out := bytes.NewBuffer([]byte{})
	if err = tmpl.Execute(out, results); err != nil {
		return "", err
	}

	return out.String(), nil
}
