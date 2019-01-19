package macTimeTracker

import (
	"bytes"
	"fmt"
	"os"

	"howett.net/plist"

	"github.com/jansorg/tom/go-tom/context"
)

func Import(filename string, ctx *context.GoTimeContext) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	decoder := plist.NewDecoder(file)
	data := map[string]interface{}{}
	err = decoder.Decode(&data)

	if err == nil {
		for _, v := range data {
			// fmt.Printf("%v\n", k)

			if m, ok := v.([]interface{}); ok {
				data := m[18]
				if d, ok := data.(map[string]interface{}); ok {
					entries, ok := d["NS.data"]
					if ok {
						if b, ok := entries.([]byte); ok {
							dec := plist.NewDecoder(bytes.NewReader(b))
							list := map[string]interface{}{}
							err = dec.Decode(&list)
							if err != nil {
								fmt.Println("error reading data bytes")
							} else {
								for k1, v1 := range list {
									fmt.Printf("%s = %v\n\n", k1, v1)
								}
							}
						}
					}
				}
			}
		}
	} else {
		fmt.Printf("error: %v", err)
	}

	return err
}
