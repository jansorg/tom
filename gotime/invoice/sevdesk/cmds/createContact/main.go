package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jansorg/gotime/gotime/invoice/sevdesk"
)

func main() {
	c := sevdesk.NewClient(os.Args[1])
	c.Logging = true
	if err := c.LoadBasicData(); err != nil {
		log.Fatal(err)
	}

	companyContact := c.NewCompanyContact("New customer", "My new customer is this")
	contact, err := c.CreateCompanyContact(companyContact)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully created contact: %s\n", contact.ID)
}
