package htmlreport

import (
	"bytes"
	"html/template"
	"path/filepath"
	"time"

	"github.com/jansorg/gotime/gotime/context"
	"github.com/jansorg/gotime/gotime/report"
)

type Report struct {
	templatePath string
	ctx          *context.GoTimeContext
}

func NewReport(templatePath string, ctx *context.GoTimeContext) *Report {
	return &Report{
		templatePath: templatePath,
		ctx:          ctx,
	}
}

func (r *Report) Render(results *report.BucketReport) (string, error) {
	tmpl, err := template.New(filepath.Base(r.templatePath)).Funcs(map[string]interface{}{
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
			return r.ctx.DurationPrinter.Minimal(duration)
		},
		"shortDuration": func(duration time.Duration) string {
			return r.ctx.DurationPrinter.Short(duration)
		},
		"longDuration": func(duration time.Duration) string {
			return r.ctx.DurationPrinter.Long(duration)
		},
	}).ParseFiles(r.templatePath, filepath.Join(filepath.Dir(r.templatePath), "commons.gohtml"))

	if err != nil {
		return "", err
	}

	out := bytes.NewBuffer([]byte{})
	if err = tmpl.Execute(out, results); err != nil {
		return "", err
	}

	return out.String(), nil
}
