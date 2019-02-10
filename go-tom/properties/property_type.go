package properties

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/rhymond/go-money"
)

var ErrUnsupportedValue = fmt.Errorf("unsupported value type")
var ErrUnsupportedUnit = fmt.Errorf("unsupported currency unit")

type PropertyType interface {
	ID() string
	Marshall(value []byte) (PropertyValue, error)
	Unmarshall(value PropertyValue) ([]byte, error)
	Parse(value string, propertyID string) (PropertyValue, error)
}

type currencyType struct {
}

func (c *currencyType) ID() string {
	return "currency"
}

func (currencyType) Marshall(value []byte) (PropertyValue, error) {
	propValue := &CurrencyValue{}
	if err := json.Unmarshal(value, &propValue); err != nil {
		return nil, err
	}
	return propValue, nil
}

func (c *currencyType) Unmarshall(value PropertyValue) ([]byte, error) {
	if value.Type().ID() != c.ID() {
		return nil, ErrUnsupportedValue
	}

	return json.Marshal(value)
}

func (c *currencyType) Parse(value string, propertyID string) (PropertyValue, error) {
	parts := strings.Split(value, " ")
	if len(parts) != 2 {
		return nil, ErrUnsupportedValue
	}

	var unit string
	var amountString string

	if isNumber(parts[0]) {
		amountString = parts[0]
		unit = parts[1]
	} else {
		unit = parts[0]
		amountString = parts[1]
	}

	currency := money.GetCurrency(unit)
	if currency == nil {
		return nil, ErrUnsupportedUnit
	}

	amount, err := parseAmount(amountString, currency)
	if err != nil {
		return nil, err
	}

	return &CurrencyValue{
		value:      amount,
		propertyID: propertyID,
	}, nil
}

func parseAmount(v string, currency *money.Currency) (*money.Money, error) {
	var parts []string
	if strings.Contains(v, ".") {
		parts = strings.Split(v, ".")
	} else if strings.Contains(v, ",") {
		parts = strings.Split(v, ",")
	} else if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		parts = []string{v, "0"}
	} else {
		return nil, ErrUnsupportedValue
	}

	if len(parts) != 2 {
		return nil, ErrUnsupportedValue
	}

	var ints []int64
	for _, p := range parts {
		if i, err := strconv.ParseInt(p, 10, 64); err != nil {
			return nil, err
		} else {
			ints = append(ints, i)
		}
	}

	var result int64
	result = int64(float64(ints[0]) * math.Pow10(currency.Fraction))
	result += int64(float64(ints[1]) * math.Pow10(currency.Fraction) / math.Pow10(currency.Fraction))

	return money.New(result, currency.Code), nil
}

func isNumber(v string) bool {
	first := string(v[0])
	var err error
	if _, err = strconv.Atoi(first); err != nil {
		_, err = strconv.ParseFloat(first, 32)
	}
	return err == nil
}
