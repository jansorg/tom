package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jansorg/gotime/gotime/invoice/sevdesk"
)

func main() {
	c := sevdesk.NewClient(os.Args[1])

	invoice, err := sevdesk.NewInvoice(sevdesk.TypeInvoice, time.Now())
	invoice.Header = "Rechnung NEU"
	invoice.Contact.ID = "7067576"
	invoice.Contact.ObjectName = "Contact"
	invoice.ContactPerson.ID = "254513"
	invoice.ContactPerson.ObjectName = "SevUser"
	invoice.Status = 100
	invoice.TaxRate = "19"
	invoice.TaxText = ""
	invoice.TaxType = sevdesk.TaxTypeNotEU
	invoice.Currency = sevdesk.USD
	invoice.DiscountTime = 0
	invoice.Address = "Joachim Ansorg\nHi"

	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.CreateInvoice(invoice)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created invoice")
	// bytes, _ := json.MarshalIndent(resp, "", "  ")
	// fmt.Printf(string(bytes))

	fmt.Println("Adding position")

	pos, err := sevdesk.NewInvoicePosition(resp.ID, "PyCharm development", c.GetQuantity(32.0, "hours"), 71.0, 0)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.CreateInvoicePos(pos)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.BrowserURL())
}
