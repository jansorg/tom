package html

import (
	"bytes"
	"html/template"

	".."
)

type Report struct {
	templatePath string
}

func NewReport(templatePath string) *Report {
	return &Report{
		templatePath: templatePath,
	}
}

func (r *Report) Render(results report.Results) (string, error) {
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
