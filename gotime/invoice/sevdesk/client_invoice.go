package sevdesk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (api *Client) CreateInvoice(data Invoice) (*InvoiceResponse, error) {
	if data.InvoiceID == "" {
		if id, err := api.FetchNextInvoiceID(data.InvoiceType, true); err != nil {
			return nil, err
		} else {
			data.InvoiceID = id
		}
	}

	req, err := api.newFormRequest("POST", "/Invoice", nil, data.asFormEncoded())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	resp, err = api.do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code %s", resp.Status)
	}

	invoiceResp := &InvoiceResponse{}
	if err = api.unwrapJSONResponse(resp, invoiceResp); err != nil {
		return nil, err
	}
	return invoiceResp, nil
}

func (api *Client) CreateInvoicePos(data InvoicePosition) (*InvoicePosResponse, error) {
	req, err := api.newFormRequest("POST", "/InvoicePos", nil, data.asFormEncoded())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	resp, err = api.do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code %s", resp.Status)
	}

	invoiceResp := &InvoicePosResponse{}
	if err = api.unwrapJSONResponse(resp, invoiceResp); err != nil {
		return nil, err
	}
	return invoiceResp, nil
}

func (api *Client) FetchNextInvoiceID(invoiceType InvoiceType, nextID bool) (string, error) {
	var err error
	var req *http.Request
	if req, err = api.newRequest("GET", "/Invoice/Factory/getNextInvoiceNumber", nil, map[string]string{
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
