package sevdesk

import (
	"fmt"
	"net/http"
	"time"
)

func (api *Client) NewInvoice(invoiceType InvoiceType, invoiceDate time.Time, header string, contactID string, status int, taxRate float32, taxText string, taxType TaxType, currency Currency, discountTime int, address string) (Invoice, error) {
	if invoiceType == "" {
		return Invoice{}, fmt.Errorf("invoiceType is empty")
	}
	if invoiceDate.IsZero() {
		return Invoice{}, fmt.Errorf("invoiceDate is required")
	}
	if header == "" {
		return Invoice{}, fmt.Errorf("header is required")
	}

	return Invoice{
		InvoiceType:   invoiceType,
		InvoiceDate:   invoiceDate,
		Header:        header,
		Contact:       IDWithType{ID: contactID, ObjectName: "Contact"},
		ContactPerson: api.userContactPerson,
		Status:        status,
		TaxRate:       taxRate,
		TaxText:       taxText,
		TaxType:       taxType,
		Currency:      currency,
		DiscountTime:  discountTime,
		Address:       address,
	}, nil
}

func (api *Client) NewQuantity(quantity float32, unitName string) (Quantity, error) {
	if unit, err := api.FindUnit(unitName); err != nil {
		return Quantity{}, err
	} else {
		return Quantity{Quantity: quantity, UnitID: unit.ID}, nil
	}
}

func (api *Client) NewInvoicePosition(invoiceID string, name string, quantity float32, unitName string, price float32, taxRate int) (InvoicePosition, error) {
	q, err := api.NewQuantity(quantity, unitName)
	if err != nil {
		return InvoicePosition{}, err
	}

	return InvoicePosition{
		InvoiceID: invoiceID,
		Name:      name,
		Quantity:  q,
		Price:     price,
		TaxRate:   taxRate,
	}, nil
}

func (api *Client) GetInvoices() (*[]InvoiceResponse, error) {
	resp, err := api.doRequest("GET", "/Invoice", nil, nil)
	if err != nil {
		return nil, err
	}

	var invoices []InvoiceResponse
	err = api.unwrapJSONResponse(resp, &invoices)
	if err != nil {
		return nil, err
	}
	return &invoices, nil
}

func (api *Client) DeleteInvoice(id string) error {
	// fixme escape ID in path?
	resp, err := api.doRequest("DELETE", fmt.Sprintf("/Invoice/%s", id), nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %s", resp.Status)
	}
	return nil
}

func (api *Client) CreateInvoice(invoicdDef Invoice) (*InvoiceResponse, error) {
	if invoicdDef.InvoiceID == "" {
		if id, err := api.FetchNextInvoiceID(invoicdDef.InvoiceType, true); err != nil {
			return nil, err
		} else {
			invoicdDef.InvoiceID = id
		}
	}

	resp, err := api.doFormRequest("POST", "/Invoice", nil, invoicdDef.asFormEncoded())
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code %s", resp.Status)
	}

	invoiceResp := &InvoiceResponse{}
	if err = api.unwrapJSONResponse(resp, invoiceResp); err != nil {
		return nil, err
	}
	return invoiceResp, nil
}

func (api *Client) CreateInvoicePos(data InvoicePosition) (*InvoicePosResponse, error) {
	resp, err := api.doFormRequest("POST", "/InvoicePos", nil, data.asFormEncoded())
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
	resp, err := api.doRequest("GET", "/Invoice/Factory/getNextInvoiceNumber", nil, map[string]string{
		"invoiceType":   string(invoiceType),
		"useNextNumber": boolToString(nextID),
	})
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	var id string
	if err := api.unwrapJSONResponse(resp, &id); err != nil {
		return "", err
	}
	return id, nil
}
