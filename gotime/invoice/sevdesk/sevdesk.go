package sevdesk

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: "https://my.sevdesk.de/api/v1",
		apiKey:  apiKey,
		http:    &http.Client{},
	}
}

type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func (api *Client) CreateInvoice(data Invoice) (string, error) {
	if data.InvoiceID == "" {
		if id, err := api.FetchNextInvoiceID(data.InvoiceType, true); err != nil {
			return "", err
		} else {
			data.InvoiceID = id
		}
	}

	req, err := api.newFormUrlencodedRequest("POST", "/Invoice", nil, data.asFormEncoded())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return "", err
	}

	var resp *http.Response
	resp, err = api.do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code %s", resp.Status)
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	return string(bytes), err
}

func iso8601(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func (api *Client) FetchNextInvoiceID(invoiceType InvoiceType, nextID bool) (string, error) {
	var err error
	var req *http.Request
	if req, err = api.newRequest("GET", "/Invoice/Factory/getNextInvoiceNumber", true, nil, map[string]string{
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

func (api *Client) newRequest(method string, path string, addToken bool, body io.Reader, query map[string]string) (*http.Request, error) {
	var u *url.URL
	var err error
	if u, err = api.makeURL(path); err != nil {
		return nil, err
	}

	q := u.Query()
	if addToken {
		q.Add("token", api.apiKey)
	}
	for k, v := range query {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	urlString := u.String()
	var req *http.Request
	if req, err = http.NewRequest(method, urlString, body); err != nil {
		return nil, err
	}

	if addToken {
		req.Header.Add("Authorization", api.apiKey)
	}
	return req, err
}

func (api *Client) newFormUrlencodedRequest(method string, path string, query map[string]string, body map[string]string) (*http.Request, error) {
	req, err := api.newRequest(method, path, false, strings.NewReader(api.createFormValues(body, true).Encode()), query)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", api.apiKey)
	return req, err
}

func (api *Client) createFormValues(data map[string]string, skipEmpty bool) url.Values {
	u := url.Values{}
	for k, v := range data {
		if !skipEmpty || v != "" {
			u.Add(k, v)
		}
	}
	return u
}

func (api *Client) do(req *http.Request) (*http.Response, error) {
	log.Printf("%s %s", req.Method, req.URL.String())
	return api.http.Do(req)
}

func (api *Client) makeURL(path string) (*url.URL, error) {
	u := fmt.Sprintf("%s%s", api.baseURL, path)
	return url.Parse(u)
}

func boolToString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
