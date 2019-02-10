package properties

import (
	"encoding/json"
	"fmt"
)

type PropertyValues []PropertyValue

type typedValue struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"value"`
}

func (p *PropertyValues) UnmarshalJSON(data []byte) error {
	var mapped []typedValue
	if err := json.Unmarshal(data, &mapped); err != nil {
		return err
	}

	for _, item := range mapped {
		if item.Type == "c" {
			var value CurrencyValue
			if err := json.Unmarshal(item.Data, &value); err != nil {
				return err
			}
			*p = append(*p, &value)
		} else {
			return fmt.Errorf("unsupported property type '%s'", item.Type)
		}

	}
	return nil
}

func (p *PropertyValues) MarshalJSON() ([]byte, error) {
	var mapped []typedValue
	for _, item := range *p {
		if _, ok := item.(*CurrencyValue); ok {
			itemData, err := json.Marshal(item)
			if err != nil {
				return nil, err
			}

			mapped = append(mapped, typedValue{
				Type: "c",
				Data: itemData,
			})
		}
	}

	return json.Marshal(mapped)
}
