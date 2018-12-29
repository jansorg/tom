package htmlreport

import (
	"bytes"
	"html/template"
	"path/filepath"
	"time"

	"github.com/jansorg/gotime/gotime/report"
)

type Report struct {
	templatePath string
}

func NewReport(templatePath string) *Report {
	return &Report{
		templatePath: templatePath,
	}
}

func (r *Report) Render(results *report.BucketReport) (string, error) {
	tmpl, err := template.New(filepath.Base(r.templatePath)).Funcs(map[string]interface{}{
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
