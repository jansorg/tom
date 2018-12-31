package converter

import "io"

type PDFConverter interface {
	ConvertHTML(in io.Reader, out io.Writer) error
}
