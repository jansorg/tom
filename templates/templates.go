// Package templates provides the built-in report templates, embedded into the binary.
package templates

import "embed"

// FS contains the built-in report templates, e.g. "reports/html/default.gohtml".
//
//go:embed reports
var FS embed.FS
