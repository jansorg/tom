package sevdesk

import (
	"fmt"
	"strconv"
	"time"
)

type InvoiceType string

const TypeInvoice InvoiceType = "RE"

type Currency string

const USD Currency = "USD"

type TaxType string

const TaxTypeNotEU TaxType = "noteu"
const TaxTypeDefault TaxType = "default"

type Person struct {
	ID   string
	Name string
}

func NewInvoice(invoiceType InvoiceType, invoiceDate time.Time) (Invoice, error) {
	if invoiceType == "" {
		return Invoice{}, fmt.Errorf("invoiceType is empty")
	}

	if invoiceDate.IsZero() {
		return Invoice{}, fmt.Errorf("invoiceDate is required")
	}

	return Invoice{
		InvoiceType: invoiceType,
		InvoiceDate: invoiceDate,
	}, nil
}

type Invoice struct {
	Header          string
	InvoiceID       string
	InvoiceType     InvoiceType
	InvoiceDate     time.Time
	DeliveryDate    time.Time
	TaxRate         string
	TaxText         string
	TaxType         TaxType
	Currency        Currency
	Discount        int
	DiscountTime    int
	Status          int
	Address         string
	Contact         Person
	ContactPerson   Person
	SmallSettlement int
}

func (data Invoice) asFormEncoded() map[string]string {
	return map[string]string{
		"header":                    data.Header,
		"invoiceNumber":             data.InvoiceID,
		"invoiceType":               string(data.InvoiceType),
		"currency":                  string(data.Currency),
		"invoiceDate":               iso8601(data.InvoiceDate),
		"deliveryDate":              iso8601(data.DeliveryDate),
		"contact[id]":               data.Contact.ID,
		"contact[objectName]":       data.Contact.Name,
		"contactPerson[id]":         data.ContactPerson.ID,
		"contactPerson[objectName]": data.ContactPerson.Name,
		"discountTime":              strconv.Itoa(data.DiscountTime),
		"taxRate":                   data.TaxRate,
		"taxText":                   stringOr0(data.TaxText),
		"taxType":                   string(data.TaxType),
		"discount":                  strconv.Itoa(data.Discount),
		"status":                    strconv.Itoa(data.Status),
		"address":                   data.Address,
	}
}

func stringOr0(v string) string {
	if v == "" {
		return "0"
	}
	return v
}
