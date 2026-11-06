package money

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Rhymond/go-money"
)

func Parse(value string) (*Money, error) {
	var amount float64
	var code string
	parts := strings.Fields(value)
	if len(parts) != 2 {
		return nil, fmt.Errorf("unable to parse money value %s, expected an amount and a currency code, e.g. '100 EUR'", value)
	}

	var err error
	if amount, err = parseNumber(parts[0]); err == nil {
		code = parts[1]
	} else if amount, err = parseNumber(parts[1]); err == nil {
		code = parts[0]
	} else {
		return nil, fmt.Errorf("unable to parse money value %s", value)
	}

	curr := money.GetCurrency(code)
	if curr == nil {
		return nil, fmt.Errorf("unknown currency %s", code)
	}
	return NewMoney(int64(math.Round(amount*math.Pow10(curr.Fraction))), code), nil
}

func parseNumber(s string) (float64, error) {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, nil
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return float64(i), nil
}

type Money struct {
	v *money.Money
}

func (m *Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount int64  `json:"amount"`
		Code   string `json:"code"`
	}{
		Amount: m.Amount(),
		Code:   m.CurrencyCode(),
	})
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var value = struct {
		Amount int64  `json:"amount"`
		Code   string `json:"code"`
	}{}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	m.v = money.New(value.Amount, value.Code)
	return nil
}

func NewMoney(amount int64, code string) *Money {
	return &Money{
		v: money.New(amount, code),
	}
}

func (m *Money) Add(v *Money) error {
	newV, err := m.v.Add(v.v)
	if err != nil {
		return err
	}

	m.v = newV
	return nil
}

func (m *Money) Multiple(factor float64) *Money {
	result := int64(math.Round(float64(m.Amount()) * factor))
	return NewMoney(result, m.CurrencyCode())
}

func (m *Money) Amount() int64 {
	return m.v.Amount()
}

func (m *Money) String() string {
	return m.v.Display()
}

func (m *Money) ParsableString() string {
	return fmt.Sprintf("%.2f %s", float64(m.Amount())/math.Pow10(m.Currency().Fraction), m.CurrencyCode())
}

func (m *Money) Currency() *money.Currency {
	return m.v.Currency()
}

func (m *Money) CurrencyCode() string {
	return m.v.Currency().Code
}
