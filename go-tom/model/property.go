package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var ErrUnknownPropertyType = fmt.Errorf("unknown proprety type")
var ErrUnsupportedValue = fmt.Errorf("unable to convert value")

type PropertyType int8

func (p *PropertyType) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}

	v, err := TypeFromString(name)
	if err != nil {
		return err
	}
	*p = v
	return nil
}

func TypeFromString(value string) (PropertyType, error) {
	switch value {
	case "string":
		return StringType, nil
	case "number":
		return NumericType, nil
	}
	return 0, ErrUnknownPropertyType
}

func (p PropertyType) MarshalJSON() ([]byte, error) {
	var name string
	switch (p) {
	case StringType:
		name = "string"
	case NumericType:
		name = "number"
	default:
		name = ""
	}

	return json.Marshal(name)
}

const (
	StringType PropertyType = iota + 1
	NumericType
)

type Property struct {
	ID                 string       `json:"id"`
	Name               string       `json:"name"`
	Prefix             string       `json:"prefix"`
	Suffix             string       `json:"suffix"`
	Type               PropertyType `json:"type"`
	ApplyToSubprojects bool         `json:"apply_to_subprojects"`
}

func (p *Property) Validate(value interface{}) error {
	ok := false
	switch p.Type {
	case StringType:
		_, ok = value.(string)
	case NumericType:
		_, ok = value.(int)
		if !ok {
			_, ok = value.(float32)
		}
		if !ok {
			_, ok = value.(float64)
		}
	}

	if !ok {
		return fmt.Errorf("invalid format")
	}
	return nil
}

func (p *Property) FromString(value string) (interface{}, error) {
	switch p.Type {
	case StringType:
		return value, nil
	case NumericType:
		return strconv.ParseFloat(value, 64)
	default:
		return nil, fmt.Errorf("unknown type %v", p.Type)
	}
}

func (p *Property) ToFloat(value interface{}) (float64, error) {
	switch p.Type {
	case NumericType:
		if f, ok := value.(float64); ok {
			return f, nil
		}
		if f, ok := value.(float32); ok {
			return float64(f), nil
		}
		if f, ok := value.(int); ok {
			return float64(f), nil
		}
		return 0, ErrUnsupportedValue
	default:
		return 0, ErrUnsupportedValue
	}
}
