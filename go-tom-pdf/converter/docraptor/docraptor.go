package docraptor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jansorg/tom/go-tom-pdf/converter"
)

var httpClient = &http.Client{Timeout: 2 * time.Minute}

const baseURL = "https://docraptor.com/docs"

func NewConverter(key string, requestConfig Request) converter.PDFConverter {
	return &docraptor{
		baseURL:       baseURL,
		key:           key,
		requestConfig: requestConfig,
	}
}

type docraptor struct {
	baseURL       string
	key           string
	requestConfig Request
}

func (a *docraptor) ConvertHTML(in io.Reader, out io.Writer) error {
	body, err := json.Marshal(a.requestConfig)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", a.baseURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error sending HTTP POST: %v", err)
	}
	req.SetBasicAuth(a.key, "")
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP POST: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("error writing into output PDF file: %s", err.Error())
	}
	return nil
}

type DocumentType string

const PDF DocumentType = "pdf"
const XLS DocumentType = "xls"
const XLSX DocumentType = "xlsx"

type MediaType string

const PRINT MediaType = "print"
const SCREEN MediaType = "screen"

type PrinceOptions struct {
	Media   MediaType `json:"media,omitempty"`
	BaseURL string    `json:"baseurl,omitempty"`
}

type Request struct {
	Test            bool          `json:"test"`
	Name            string        `json:"name,omitempty"`
	DocumentURL     string        `json:"document_url,omitempty"`
	DocumentContent string        `json:"document_content,omitempty"`
	Type            DocumentType  `json:"type"`
	JavaScript      bool          `json:"javascript,omitempty"`
	Pipeline        int8          `json:"pipeline,omitempty"`
	PrinceOptions   PrinceOptions `json:"prince_options,omitempty"`
}

type Response struct {
	PDFUrl  string  `json:"pdf"`
	MBIn    float32 `json:"mbIn"`
	MBOut   float32 `json:"mbOut"`
	Cost    float32 `json:"cost"`
	Success bool    `json:"success"`
}
