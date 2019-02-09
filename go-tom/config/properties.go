package config
//
// import (
// 	"fmt"
// 	"strconv"
//
// 	"github.com/jansorg/tom/go-tom/model"
// )
//
// type Property int
// type StringProperty Property
// type FloatProperty Property
// type IntProperty Property
//
// // predefined list of properties
// const (
// 	ContactIDProperty          StringProperty = iota
// 	InvoiceHourlyRateProperty  FloatProperty  = iota
// 	InvoiceTaxRateProperty     FloatProperty  = iota
// 	InvoiceLanguageProperty    StringProperty = iota
// 	InvoiceDescriptionProperty StringProperty = iota
// 	InvoiceCurrencyProperty    StringProperty = iota
// 	InvoiceAddressProperty     StringProperty = iota
// )
//
// func (p StringProperty) key() string {
// 	switch p {
// 	case ContactIDProperty:
// 		return "contactID"
// 	case InvoiceDescriptionProperty:
// 		return "invoiceDescription"
// 	case InvoiceLanguageProperty:
// 		return "invoiceLang"
// 	case InvoiceCurrencyProperty:
// 		return "invoiceCurrency"
// 	case InvoiceAddressProperty:
// 		return "invoiceAddress"
// 	default:
// 		return ""
// 	}
// }
//
// func (p StringProperty) Get(source model.PropertyHolder) (string, bool) {
// 	name := p.key()
// 	v, ok := source.GetProperties()[name]
// 	if !ok {
// 		return "", false
// 	}
// 	return v, ok
// }
//
// func (p StringProperty) Set(value string, target model.PropertyHolder) {
// 	name := p.key()
// 	target.GetProperties()[name] = value
// }
//
// func (p FloatProperty) key() string {
// 	switch p {
// 	case InvoiceHourlyRateProperty:
// 		return "hourlyRate"
// 	case InvoiceTaxRateProperty:
// 		return "taxRate"
// 	default:
// 		return ""
// 	}
// }
//
// func (p FloatProperty) Get(source model.PropertyHolder) (float64, bool) {
// 	name := p.key()
// 	v, ok := source.GetProperties()[name]
// 	if !ok {
// 		return 0, false
// 	}
//
// 	f, err := strconv.ParseFloat(v, 64)
// 	if err != nil {
// 		return 0, false
// 	}
//
// 	return float64(f), true
// }
//
// func (p FloatProperty) Set(value float64, target model.PropertyHolder) {
// 	name := p.key()
// 	target.GetProperties()[name] = fmt.Sprintf("%.4f", value)
// }
//
// func (p IntProperty) key() string {
// 	switch p {
// 	default:
// 		return ""
// 	}
// }
//
// func (p IntProperty) Get(source model.PropertyHolder) (int64, bool) {
// 	name := p.key()
// 	v, ok := source.GetProperties()[name]
// 	if !ok {
// 		return 0, false
// 	}
//
// 	i, err := strconv.ParseInt(v, 10, 64)
// 	if err != nil {
// 		return 0, false
// 	}
// 	return i, true
// }
//
// func (p IntProperty) Set(value int64, target model.PropertyHolder) {
// 	name := p.key()
// 	target.GetProperties()[name] = fmt.Sprintf("%d", value)
// }
