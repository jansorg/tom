// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p := messageKeyToIndex[key]
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"de": &dictionary{index: deIndex, data: deData},
		"en": &dictionary{index: enIndex, data: enData},
	}
	fallback := language.MustParse("en")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"%.2f":                0,
	"Date":                11,
	"Duration":            2,
	"End":                 13,
	"Exact duration":      4,
	"Exact sales":         3,
	"Exact tracked time:": 10,
	"Notes":               14,
	"Project":             7,
	"Rounded duration":    6,
	"Sales":               1,
	"Start":               12,
	"Time range:":         8,
	"Total":               5,
	"Tracked time:":       9,
}

var deIndex = []uint32{ // 16 elements
	0x00000000, 0x00000008, 0x0000000f, 0x00000015,
	0x00000024, 0x00000031, 0x00000038, 0x00000047,
	0x0000004f, 0x0000005c, 0x0000006b, 0x00000080,
	0x00000086, 0x0000008d, 0x00000092, 0x0000009e,
} // Size: 88 bytes

const deData string = "" + // Size: 158 bytes
	"\x02%.2[1]f\x02Umsatz\x02Dauer\x02Exakter Umsatz\x02Exakte Dauer\x02Gesa" +
	"mt\x02Gerundete Zeit\x02Projekt\x02Zeitbereich:\x02Erfasste Zeit:\x02Exa" +
	"kt erfasste Zeit:\x02Datum\x02Beginn\x02Ende\x02Anmerkungen"

var enIndex = []uint32{ // 16 elements
	0x00000000, 0x00000008, 0x0000000e, 0x00000017,
	0x00000023, 0x00000032, 0x00000038, 0x00000049,
	0x00000051, 0x0000005d, 0x0000006b, 0x0000007f,
	0x00000084, 0x0000008a, 0x0000008e, 0x00000094,
} // Size: 88 bytes

const enData string = "" + // Size: 148 bytes
	"\x02%.2[1]f\x02Sales\x02Duration\x02Exact sales\x02Exact duration\x02Tot" +
	"al\x02Rounded duration\x02Project\x02Time range:\x02Tracked time:\x02Exa" +
	"ct tracked time:\x02Date\x02Start\x02End\x02Notes"

	// Total table size 482 bytes (0KiB); checksum: 774B414F
