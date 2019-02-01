package htmlreport

import (
	"bytes"
	htmlTemplate "html/template"
	"path"
	"path/filepath"
	"time"

	assetTemplate "github.com/arschles/go-bindata-html-template"

	"github.com/jansorg/tom/go-tom"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/dateUtil"
	"github.com/jansorg/tom/go-tom/report"
)

type Report struct {
	workingDir string
	options    Options
	ctx        *context.TomContext
}

type Options struct {
	DecimalDuration  bool
	TemplateName     string
	TemplateFilePath string
}

func NewReport(workingDir string, opts Options, ctx *context.TomContext) *Report {
	return &Report{
		options:    opts,
		workingDir: workingDir,
		ctx:        ctx,
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
			if r.options.DecimalDuration {
				return r.ctx.DecimalDurationPrinter.Minimal(duration)
			}
			return r.ctx.DurationPrinter.Minimal(duration)
		},
		"shortDuration": func(duration time.Duration) string {
			if r.options.DecimalDuration {
				return r.ctx.DecimalDurationPrinter.Short(duration)
			}
			return r.ctx.DurationPrinter.Short(duration)
		},
		"longDuration": func(duration time.Duration) string {
			if r.options.DecimalDuration {
				return r.ctx.DecimalDurationPrinter.Long(duration)
			}
			return r.ctx.DurationPrinter.Long(duration)
		},
		"hasFlag": func(name string) bool {
			if name == "showSummary" {
				return true
			}
			return false
		},
		"isMatrix": func(bucket report.ResultBucket) bool {
			if bucket.Depth() != 2 {
				return false
			}

			// all buckets must have the same number of children
			refCol := bucket.ChildBuckets[0].ChildBuckets
			for _, b := range bucket.ChildBuckets {
				if len(b.ChildBuckets) != len(refCol) {
					return false
				}

				for i, col := range b.ChildBuckets {
					if col.Title() != refCol[i].Title() {
						// fmt.Printf("title: %s <> %s", col.Title(), refCol[i].Title())
						return false
					}
				}
			}

			return true
		},
		"sumChildValues": func(parent report.ResultBucket, childIndex int) *dateUtil.DurationSum {
			sum := dateUtil.NewDurationSum()

			for _, b := range parent.ChildBuckets {
				if len(b.ChildBuckets) >= childIndex {
					sum.AddSum(b.ChildBuckets[childIndex].Duration)
				}
			}

			return sum
		},
	}

	if r.options.TemplateFilePath != "" {
		templatePath := r.options.TemplateFilePath
		if !filepath.IsAbs(templatePath) {
			templatePath = filepath.Join(r.workingDir, templatePath)
		}

		files, err := filepath.Glob(filepath.Join(filepath.Dir(templatePath), "*.gohtml"))

		tmpl, err := htmlTemplate.New(filepath.Base(templatePath)).Funcs(functionMap).ParseFiles(append(files, templatePath)...)
		if err != nil {
			return "", err
		}
		out := bytes.NewBuffer([]byte{})
		if err = tmpl.Execute(out, results); err != nil {
			return "", err
		}

		return out.String(), nil
	} else {
		templatePath := r.options.TemplateName

		baseDir := path.Join("reports", "html")
		templateFiles := []string{
			path.Join(baseDir, templatePath+".gohtml"),
			path.Join(baseDir, "commons.gohtml"),
		}

		tmpl, err := assetTemplate.New(templatePath, tom.Asset).Funcs(functionMap).ParseFiles(templateFiles...)
		if err != nil {
			return "", err
		}
		out := bytes.NewBuffer([]byte{})
		if err = tmpl.Execute(out, results); err != nil {
			return "", err
		}

		return out.String(), nil
	}
}
