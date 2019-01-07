package sevdesk

import (
	"net/http"
	"strconv"
	"time"
)

type CompanyContact struct {
	Name           string
	Status         int
	CustomerNumber string
	// parent
	CategoryID string
	// category[id]: 3
	// category[objectName]: Category
	Description string
	// defaultCashbackTime: null
	// defaultCashbackPercent: null
	// defaultTimeToPay: null
	// taxType: null
	// taxSet: null
	// defaultDiscountPercentage: true
	// objectName: Contact
	// types: [object Object]
}

func (c CompanyContact) asFormEncoded() map[string]string {
	return map[string]string{
		"name":                 c.Name,
		"description":          c.Description,
		"status":               strconv.Itoa(c.Status),
		"customerNumber":       c.CustomerNumber,
		"category[id]":         c.CategoryID,
		"category[objectName]": "Category",
	}
}

// type PersonContact struct{}

type ContactResponse struct {
	ID         string `json:"id"`
	ObjectName string `json:"objectName"`
	// AdditionalInformation     interface{} `json:"additionalInformation"`
	Create                    time.Time   `json:"create"`
	Update                    time.Time   `json:"update"`
	Name                      string      `json:"name"`
	Status                    string      `json:"status"`
	CustomerNumber            string      `json:"customerNumber"`
	Surename                  string      `json:"surename"`
	Familyname                string      `json:"familyname"`
	Title                     string      `json:"titel"`
	Category                  IDWithType  `json:"category"`
	Description               string      `json:"description"`
	AcademicTitle             string      `json:"academicTitle"`
	Gender                    string      `json:"gender"`
	SevClient                 IDWithType  `json:"sevClient"`
	Name2                     string      `json:"name2"`
	Birthday                  interface{} `json:"birthday"`
	VatNumber                 interface{} `json:"vatNumber"`
	BankAccount               interface{} `json:"bankAccount"`
	BankNumber                interface{} `json:"bankNumber"`
	DefaultCashbackTime       interface{} `json:"defaultCashbackTime"`
	DefaultCashbackPercent    interface{} `json:"defaultCashbackPercent"`
	DefaultTimeToPay          interface{} `json:"defaultTimeToPay"`
	TaxNumber                 interface{} `json:"taxNumber"`
	TaxOffice                 interface{} `json:"taxOffice"`
	ExemptVat                 string      `json:"exemptVat"`
	TaxType                   interface{} `json:"taxType"`
	DefaultDiscountAmount     interface{} `json:"defaultDiscountAmount"`
	DefaultDiscountPercentage string      `json:"defaultDiscountPercentage"`
}

func (api *Client) NewCompanyContact(name string, description string) CompanyContact {
	return CompanyContact{
		Name:        name,
		Description: description,
		CategoryID:  "3",
	}
}

func (api *Client) GetNextCustomerNumber() (string, error) {
	resp, err := api.doRequest("GET", "/Contact/Factory/getNextCustomerNumber", nil, nil)
	if err != nil {
		return "", err
	}

	id := ""
	err = api.unwrapJSONResponse(resp, &id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (api *Client) CreateCompanyContact(contact CompanyContact) (*ContactResponse, error) {
	if contact.CustomerNumber == "" {
		nextID, err := api.GetNextCustomerNumber()
		if err != nil {
			return nil, err
		}
		contact.CustomerNumber = nextID
	}

	contactDef := ContactResponse{}
	return &contactDef, api.postForm("/Contact", []int{http.StatusCreated}, contact.asFormEncoded(), &contactDef)
}

func (api *Client) GetContacts() ([]ContactResponse, error) {
	var contacts []ContactResponse
	return contacts, api.getJSON("/Contact", &contacts)
}
