package htmlreport

import (
	"bytes"
	"html/template"

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

func (r *Report) Render(results report.ResultBucket) (string, error) {
	tmpl, err := template.ParseFiles(r.templatePath)
	if err != nil {
		return "", err
	}

	out := bytes.NewBuffer([]byte{})
	if err = tmpl.Execute(out, results); err != nil {
		return "", err
	}

	return out.String(), nil
}
