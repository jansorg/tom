package sevdesk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func NewSevdeskClient(apiKey string) *SevdeskClient {
	return &SevdeskClient{
		baseURL: "https://my.sevdesk.de/api/v1",
		apiKey:  apiKey,
		http:    &http.Client{},
	}
}

type SevdeskClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

type InvoiceType string

const TypeInvoice = "RE"

func (api *SevdeskClient) FetchNextInvoiceID(invoiceType InvoiceType, nextID bool) (string, error) {
	var err error
	var req *http.Request
	if req, err = api.makeRequest("GET", "/Invoice/Factory/getNextInvoiceNumber", true, map[string]string{
		"invoiceType":   string(invoiceType),
		"useNextNumber": boolToString(nextID),
	}); err != nil {
		return "", err
	}

	var resp *http.Response
	if resp, err = api.do(req); err != nil {
		return "", err
	} else if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	idValue := struct {
		Objects string `json:"objects"`
	}{}
	if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	} else {
		if err = json.Unmarshal(bodyBytes, &idValue); err != nil {
			return "", err
		}
	}

	return idValue.Objects, nil
}

func (api *SevdeskClient) makeRequest(method string, path string, addToken bool, query map[string]string) (*http.Request, error) {
	var u *url.URL
	var err error
	if u, err = api.makeURL(path); err != nil {
		return nil, err
	}

	q := u.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	if addToken && query["token"] == "" {
		query["token"] = api.apiKey
	}
	u.RawQuery = q.Encode()

	urlString := u.String()
	var req *http.Request
	if req, err = http.NewRequest(method, urlString, nil); err != nil {
		return nil, err
	}

	if addToken {
		req.Header.Add("Authorization", api.apiKey)
	}
	return req, err
}

func (api *SevdeskClient) do(req *http.Request) (*http.Response, error) {
	log.Printf("%s %s", req.Method, req.URL.String())
	return api.http.Do(req)
}

func (api *SevdeskClient) makeURL(path string) (*url.URL, error) {
	u := fmt.Sprintf("%s%s", api.baseURL, path)
	return url.Parse(u)
}

func boolToString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
