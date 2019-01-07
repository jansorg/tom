package api2pdf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/jansorg/gotime/go-tom-pdf/converter"
)

const baseURL = "https://v2018.api2pdf.com"

type RequestType int8

const (
	RequestWkHtml RequestType = iota + 1
	RequestChrome
)

func (t RequestType) urlPath() string {
	if t == RequestWkHtml {
		return "/wkhtmltopdf/html"
	} else if t == RequestChrome {
		return "/chrome/html"
	}
	return ""
}

func NewConverter(key string, requestType RequestType, requestConfig Request) converter.PDFConverter {
	return &api2pdf{
		baseURL:       baseURL,
		key:           key,
		reuestType:    requestType,
		requestConfig: requestConfig,
	}
}

type api2pdf struct {
	baseURL       string
	key           string
	reuestType    RequestType
	requestConfig Request
}

func (a *api2pdf) ConvertHTML(in io.Reader, out io.Writer) error {
	path := a.reuestType.urlPath()
	if path == "" {
		return fmt.Errorf("unknown request type %v", a.reuestType)
	}

	if a.requestConfig.HTML == "" {
		htmlBytes, err := ioutil.ReadAll(in)
		if err != nil {
			return fmt.Errorf("error reading input file: %s", err.Error())
		}
		a.requestConfig.HTML = string(htmlBytes)
	}

	body, err := json.Marshal(a.requestConfig)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", a.baseURL, path), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error sending HTTP POST: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", a.key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP POST: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading POST response: %v", err)
	}

	var r Response
	err = json.Unmarshal(jsonBody, &r)
	if err != nil {
		return fmt.Errorf("error unmarshalling POST response: %v", err)
	}

	if !r.Success {
		return fmt.Errorf("response indicated that it wasn't successful")
	}

	pdfResp, err := http.DefaultClient.Get(r.PDFUrl)
	if err != nil {
		return fmt.Errorf("PDF url returned unexpected status code: %s", err)
	}
	if pdfResp.StatusCode != http.StatusOK {
		return fmt.Errorf("PDF url returned unexpected status code: %d", pdfResp.StatusCode)
	}

	if _, err = io.Copy(out, pdfResp.Body); err != nil {
		return fmt.Errorf("error writing into output PDF file: %s", err.Error())
	}
	return nil
}

type PageOptions struct {
	Orientation string `json:"orientation,omitempty"`
	PageSize    string `json:"pageSize,omitempty"`
}

type Request struct {
	HTML      string       `json:"html,required"`
	InlinePDF bool         `json:"inlinePdf,omitempty"`
	FileName  string       `json:"fileName,omitempty"`
	Options   *PageOptions `json:"options,omitempty"`
}

type Response struct {
	PDFUrl  string  `json:"pdf"`
	MBIn    float32 `json:"mbIn"`
	MBOut   float32 `json:"mbOut"`
	Cost    float32 `json:"cost"`
	Success bool    `json:"success"`
}
