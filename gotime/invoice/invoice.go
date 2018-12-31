package invoice

type Factory interface {
	Name() string
	CreateInvoice(req Request)
}

type Request struct {
	// 	'header' => "Invoice " . $nextInvoiceNumber,
	// 	'invoiceNumber' => $nextInvoiceNumber,
	// 	'invoiceType' => "RE",
	// 	'invoiceDate' => $invoiceDate->format(DateTime::ISO8601),
	// 	'deliveryDate' => $deliveryDate->format(DateTime::ISO8601),
	// 	'contactPerson[id]' => '11822',
	// 	'contactPerson[objectName]' => "SevUser",
	// 	'discountTime' => '0',
	// 	'taxRate' => "19",
	// 	'taxText' => 'Umsatzsteuer ausweisen',
	// 	'taxType' => "default",
	// 	'smallSettlement' => '0',
	// 	'currency' => 'EUR',
	// 	'discount' => '0',
	// 	'status' => 100,
	// 	'address' => $address
}

type Response struct {
	ID string
}
