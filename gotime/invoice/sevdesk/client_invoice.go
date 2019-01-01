package sevdesk

import (
	"fmt"
	"net/http"
	"time"
)

func (api *Client) NewInvoice(invoiceDate time.Time, header string, contactID string, status int, taxRate float64, taxText string, taxType TaxType, currency Currency, discountTime int, address string) (Invoice, error) {
	if invoiceDate.IsZero() {
		return Invoice{}, fmt.Errorf("invoiceDate is required")
	}
	if header == "" {
		return Invoice{}, fmt.Errorf("header is required")
	}

	return Invoice{
		InvoiceType:   TypeInvoice,
		InvoiceDate:   invoiceDate,
		InvoiceTitle:  header,
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

func (api *Client) NewQuantity(quantity float64, unitName string) (Quantity, error) {
	if unit, err := api.FindUnit(unitName); err != nil {
		return Quantity{}, err
	} else {
		return Quantity{Quantity: quantity, UnitID: unit.ID}, nil
	}
}

func (api *Client) NewInvoicePosition(invoiceID string, name string, quantity float64, unitName string, price float64, taxRate int) (InvoicePosition, error) {
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
	var invoices []InvoiceResponse
	return &invoices, api.getJSON("/Invoice", &invoices)
}

func (api *Client) DeleteInvoice(id string) error {
	return api.delete("Invoice", id)
}

func (api *Client) CreateInvoice(invoice Invoice) (*InvoiceResponse, error) {
	if invoice.InvoiceID == "" {
		id, err := api.FetchNextInvoiceID(invoice.InvoiceType, true)
		if err != nil {
			return nil, err
		}
		invoice.InvoiceID = id
	}

	invoiceResp := &InvoiceResponse{}
	return invoiceResp, api.postForm("/Invoice", []int{http.StatusCreated}, invoice.asFormEncoded(), invoiceResp)
}

func (api *Client) CreateInvoicePos(data InvoicePosition) (*InvoicePosResponse, error) {
	invoiceResp := &InvoicePosResponse{}
	return invoiceResp, api.postForm("/InvoicePos", []int{http.StatusCreated}, data.asFormEncoded(), invoiceResp)
}

func (api *Client) FetchNextInvoiceID(invoiceType InvoiceType, nextID bool) (string, error) {
	var id string
	return id, api.getJSON("/Invoice/Factory/getNextInvoiceNumber", &id)
}
