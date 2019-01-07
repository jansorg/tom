package createContact

import (
	"fmt"
	"log"
	"os"

	"github.com/jansorg/gotime/go-tom/invoice/sevdesk"
)

func main() {
	c := sevdesk.NewClient(os.Args[1])
	c.Logging = false
	if err := c.LoadBasicData(); err != nil {
		log.Fatal(err)
	}

	// invoice, err := c.NewInvoice(sevdesk.TypeInvoice, time.Now(), "Invoice Kite 2018-12", "7067576", 100, 0.0, "", sevdesk.TaxTypeNotEU, sevdesk.USD, 14, "")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// resp, err := c.CreateInvoice(invoice)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf("created invoice")
	// fmt.Println("Adding position")
	//
	// pos, err := c.NewInvoicePosition(resp.ID, "PyCharm development", 32.0, "hours", 71.0, 0)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// _, err = c.CreateInvoicePos(pos)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Println(resp.BrowserURL())
	//
	invoices, err := c.GetInvoices()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range *invoices {
		err = c.DeleteInvoice(i.ID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("deleted invoice %s", i.ID)
	}
}
