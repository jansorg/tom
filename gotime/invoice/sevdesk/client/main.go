package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jansorg/gotime/gotime/invoice/sevdesk"
)

func main() {
	c := sevdesk.NewSevdeskClient(os.Args[1])
	id, err := c.FetchNextInvoiceID(sevdesk.TypeInvoice, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("next ID: %s\n", id)
}
