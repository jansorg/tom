package properties

import (
	"encoding/json"

	"github.com/rhymond/go-money"
)

type CurrencyValue struct {
	propertyID string
	value      *money.Money
}

func NewCurrency(propertyID string, value *money.Money) *CurrencyValue {
	return &CurrencyValue{
		propertyID: propertyID,
		value:      value,
	}
}

func (c CurrencyValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PropertyID string `json:"property,required"`
		Value      int64  `json:"value"`
		Unit       string `json:"unit"`
	}{
		PropertyID: c.propertyID,
		Value:      c.value.Amount(),
		Unit:       c.value.Currency().Code,
	})
}

func (c *CurrencyValue) UnmarshalJSON(data []byte) error {
	value := struct {
		PropertyID string `json:"property,required"`
		Value      int64  `json:"value"`
		Unit       string `json:"unit"`
	}{}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	c.propertyID = value.PropertyID
	c.value = money.New(value.Value, value.Unit)
	return nil
}

func (c *CurrencyValue) String() string {
	return c.value.Display()
}

func (c *CurrencyValue) Type() PropertyType {
	return CurrencyType
}

func (c *CurrencyValue) IsString() bool {
	return false
}

func (c *CurrencyValue) IsFloat() bool {
	return false
}

func (c *CurrencyValue) IsCurrency() bool {
	return true
}

func (c *CurrencyValue) AsString() (string, error) {
	return "", ErrUnsupportedType
}

func (c *CurrencyValue) AsFloat() (float64, error) {
	return 0, ErrUnsupportedType
}

func (c *CurrencyValue) AsCurrency() (*money.Money, error) {
	return c.value, nil
}

func (c *CurrencyValue) PropertyID() string {
	return c.propertyID
}
