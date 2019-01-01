package config

import (
	"fmt"
	"strconv"

	"github.com/jansorg/gotime/gotime/store"
)

type Property int
type StringProperty Property
type FloatProperty Property
type IntProperty Property

// predefined list of properties
const (
	HourlyRateProperty  FloatProperty  = iota
	DescriptionProperty StringProperty = iota
	CurrencyProperty    StringProperty = iota
	AddressProperty     StringProperty = iota
)

func (p StringProperty) key() string {
	switch p {
	case DescriptionProperty:
		return "description"
	case CurrencyProperty:
		return "currency"
	case AddressProperty:
		return "address"
	default:
		return ""
	}
}

func (p StringProperty) Get(source store.PropertyHolder) (string, bool) {
	name := p.key()
	v, ok := source.GetProperties()[name]
	if !ok {
		return "", false
	}
	return v, ok
}

func (p StringProperty) Set(value string, target store.PropertyHolder) {
	name := p.key()
	target.GetProperties()[name] = value
}

func (p FloatProperty) key() string {
	switch p {
	case HourlyRateProperty:
		return "hourlyRate"
	default:
		return ""
	}
}

func (p FloatProperty) Get(source store.PropertyHolder) (float64, bool) {
	name := p.key()
	v, ok := source.GetProperties()[name]
	if !ok {
		return 0, false
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, false
	}

	return float64(f), true
}

func (p FloatProperty) Set(value float64, target store.PropertyHolder) {
	name := p.key()
	target.GetProperties()[name] = fmt.Sprintf("%.4f", value)
}

func (p IntProperty) key() string {
	switch p {
	default:
		return ""
	}
}

func (p IntProperty) Get(source store.PropertyHolder) (int64, bool) {
	name := p.key()
	v, ok := source.GetProperties()[name]
	if !ok {
		return 0, false
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}

func (p IntProperty) Set(value int64, target store.PropertyHolder) {
	name := p.key()
	target.GetProperties()[name] = fmt.Sprintf("%d", value)
}
