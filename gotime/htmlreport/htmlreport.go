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
			return r.ctx.Translator.Sprintf(key)
		},
		"formatNumber": func(n interface{}) string {
			return r.ctx.NumberFormat.Sprint(n)
		},
		"formatTime": func(date time.Time) string {
			return date.Format("15:04:05")
		},
		"formatDate": func(date time.Time) string {
			return date.Format("2006-01-02")
		},
		"formatDateTime": func(date time.Time) string {
			return date.Format("2006-01-02 15:04:05")
		},
		"formatDuration": func(duration time.Duration) string {
			return duration.String()
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
