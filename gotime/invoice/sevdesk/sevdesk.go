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

func (api *Client) newRequest(method string, path string, body io.Reader, query map[string]string) (*http.Request, error) {
	var u *url.URL
	var err error
	if u, err = api.makeURL(path); err != nil {
		return nil, err
	}

	q := u.Query()
	if method == "GET" {
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

	if method != "GET" {
		req.Header.Add("Authorization", api.apiKey)
	}
	return req, err
}

func (api *Client) newFormRequest(method string, path string, query map[string]string, body map[string]string) (*http.Request, error) {
	req, err := api.newRequest(method, path, strings.NewReader(api.createFormValues(body, true).Encode()), query)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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

func (api *Client) unwrapJSONResponse(resp *http.Response, target interface{}) error {
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	// target is a pointer and will be updated
	return json.Unmarshal(bytes, &struct {
		Objects interface{} `json:"objects"`
	}{Objects: target})
}

func (api *Client) GetQuantity(quantity float32, name string) Quantity {
	return Quantity{
		Quantity: quantity,
		UnitID:   "1",
	}
}
