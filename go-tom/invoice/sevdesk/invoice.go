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

type IDWithType struct {
	ID         string `json:"id"`
	ObjectName string `json:"objectName"`
}

type Invoice struct {
	InvoiceTitle    string
	InvoiceID       string
	InvoiceType     InvoiceType
	InvoiceDate     time.Time
	DeliveryDate    time.Time
	TaxRate         float64
	TaxText         string
	TaxType         TaxType
	Currency        Currency
	Discount        int
	DiscountTime    int
	Status          int
	Address         string
	Contact         IDWithType
	ContactPerson   IDWithType
	SmallSettlement int
}

func (data Invoice) asFormEncoded() map[string]string {
	return map[string]string{
		"header":                    data.InvoiceTitle,
		"invoiceNumber":             data.InvoiceID,
		"invoiceType":               string(data.InvoiceType),
		"currency":                  string(data.Currency),
		"invoiceDate":               iso8601(data.InvoiceDate),
		"deliveryDate":              iso8601(data.DeliveryDate),
		"contact[id]":               data.Contact.ID,
		"contact[objectName]":       data.Contact.ObjectName,
		"contactPerson[id]":         data.ContactPerson.ID,
		"contactPerson[objectName]": data.ContactPerson.ObjectName,
		"discountTime":              strconv.Itoa(data.DiscountTime),
		"taxRate":                   strconv.Itoa(int(data.TaxRate)),
		"taxText":                   stringOr0(data.TaxText),
		"taxType":                   string(data.TaxType),
		"discount":                  strconv.Itoa(data.Discount),
		"status":                    strconv.Itoa(data.Status),
		"address":                   data.Address,
	}
}

type InvoicePosition struct {
	InvoiceID string
	Name      string
	Quantity  Quantity
	Price     float64
	TaxRate   int
}

type Quantity struct {
	Quantity float64
	UnitID   string
}

func (pos InvoicePosition) asFormEncoded() map[string]string {
	return map[string]string{
		"invoice[id]":         pos.InvoiceID,
		"invoice[objectName]": "Invoice",
		"name":                pos.Name,
		"quantity":            fmt.Sprintf("%f", pos.Quantity.Quantity),
		"price":               fmt.Sprintf("%f", pos.Price),
		"unity[id]":           pos.Quantity.UnitID,
		"unity[objectName]":   "Unity",
		"taxRate":             strconv.Itoa(pos.TaxRate),
	}
}

type InvoiceResponse struct {
	ID                    string      `json:"id"`
	ObjectName            string      `json:"objectName"`
	AdditionalInformation interface{} `json:"additionalInformation"`
	InvoiceNumber         string      `json:"invoiceNumber"`
	Contact               IDWithType  `json:"contact"`
	Create                time.Time   `json:"create"`
	Update                time.Time   `json:"update"`
	InvoiceDate           time.Time   `json:"invoiceDate"`
	Header                interface{} `json:"header"`
	HeadText              string      `json:"headText"`
	FootText              string      `json:"footText"`
	TimeToPay             interface{} `json:"timeToPay"`
	DiscountTime          string      `json:"discountTime"`
	Discount              string      `json:"discount"`
	AddressName           interface{} `json:"addressName"`
	AddressStreet         interface{} `json:"addressStreet"`
	AddressZip            interface{} `json:"addressZip"`
	AddressCity           interface{} `json:"addressCity"`
	PayDate               time.Time   `json:"payDate"`
	CreateUser            IDWithType  `json:"createUser"`
	SevClient             IDWithType  `json:"sevClient"`
	DeliveryDate          interface{} `json:"deliveryDate"`
	Status                string      `json:"status"`
	SmallSettlement       string      `json:"smallSettlement"`
	ContactPerson         IDWithType  `json:"contactPerson"`
	TaxRate               string      `json:"taxRate"`
	TaxText               string      `json:"taxText"`
	DunningLevel          interface{} `json:"dunningLevel"`
	AddressParentName     interface{} `json:"addressParentName"`
	TaxType               string      `json:"taxType"`
	SendDate              interface{} `json:"sendDate"`
	OriginLastInvoice     interface{} `json:"originLastInvoice"`
	InvoiceType           string      `json:"invoiceType"`
	AccountIntervall      interface{} `json:"accountIntervall"`
	AccountLastInvoice    interface{} `json:"accountLastInvoice"`
	AccountNextInvoice    interface{} `json:"accountNextInvoice"`
	ReminderTotal         interface{} `json:"reminderTotal"`
	ReminderDebit         interface{} `json:"reminderDebit"`
	ReminderDeadline      interface{} `json:"reminderDeadline"`
	ReminderCharge        interface{} `json:"reminderCharge"`
	AddressParentName2    interface{} `json:"addressParentName2"`
	AddressName2          interface{} `json:"addressName2"`
	AddressGender         interface{} `json:"addressGender"`
	AccountEndDate        interface{} `json:"accountEndDate"`
	Address               interface{} `json:"address"`
	Currency              string      `json:"currency"`
	// SumNet                              float32     `json:"sumNet"`
	// SumTax                              float32     `json:"sumTax"`
	// SumGross                            float32     `json:"sumGross"`
	// SumDiscounts                        float32     `json:"sumDiscounts"`
	// SumNetForeignCurrency               float32     `json:"sumNetForeignCurrency"`
	// SumTaxForeignCurrency               float32     `json:"sumTaxForeignCurrency"`
	// SumGrossForeignCurrency             float32     `json:"sumGrossForeignCurrency"`
	// SumDiscountsForeignCurrency         float32     `json:"sumDiscountsForeignCurrency"`
	// SumNetAccounting                    float32     `json:"sumNetAccounting"`
	// SumTaxAccounting                    float32     `json:"sumTaxAccounting"`
	// SumGrossAccounting                  float32     `json:"sumGrossAccounting"`
	CustomerInternalNote                string      `json:"customerInternalNote"`
	ShowNet                             string      `json:"showNet"`
	Enshrined                           interface{} `json:"enshrined"`
	SendType                            string      `json:"sendType"`
	DeliveryDateUntil                   interface{} `json:"deliveryDateUntil"`
	SendPaymentReceivedNotificationDate interface{} `json:"sendPaymentReceivedNotificationDate"`
}

func (r *InvoiceResponse) BrowserURL() string {
	return fmt.Sprintf("https://my.sevdesk.de/#/fi/edit/type/%s/id/%s", r.InvoiceType, r.ID)
}

type InvoicePosResponse struct {
	ID                    string          `json:"id"`
	ObjectName            string          `json:"objectName"`
	AdditionalInformation interface{}     `json:"additionalInformation"`
	Create                time.Time       `json:"create"`
	Update                time.Time       `json:"update"`
	Invoice               InvoiceResponse `json:"invoice"`
	Quantity              string          `json:"quantity"`
	Price                 string          `json:"price"`
	Name                  string          `json:"name"`
	Priority              string          `json:"priority"`
	Unity                 IDWithType      `json:"unity"`
	SevClient             IDWithType      `json:"sevClient"`
	PositionNumber        string          `json:"positionNumber"`
	Text                  interface{}     `json:"text"`
	Discount              interface{}     `json:"discount"`
	TaxRate               string          `json:"taxRate"`
	Temporary             string          `json:"temporary"`
	SumNet                string          `json:"sumNet"`
	SumGross              string          `json:"sumGross"`
	SumDiscount           string          `json:"sumDiscount"`
	SumTax                string          `json:"sumTax"`
	SumNetAccounting      string          `json:"sumNetAccounting"`
	SumTaxAccounting      string          `json:"sumTaxAccounting"`
	SumGrossAccounting    string          `json:"sumGrossAccounting"`
	PriceNet              string          `json:"priceNet"`
	PriceGross            string          `json:"priceGross"`
	PriceTax              string          `json:"priceTax"`
}
