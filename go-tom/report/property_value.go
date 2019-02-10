package report

import (
	"fmt"
	"time"

	"github.com/rhymond/go-money"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/properties"
	"github.com/jansorg/tom/go-tom/util"
)

func NewPropertyValueSum(property *properties.Property, frameRounding util.RoundingConfig) PropertyValueSum {
	if property.TypeID == properties.CurrencyType.ID() {
		return &currencyPropertySum{
			property: property,
			rounding: frameRounding,
			exact:    make(map[string]map[int64]time.Duration),
			rounded:  make(map[string]map[int64]time.Duration),
		}
	}
	return nil
}

type PropertyValueSum interface {
	Property() *properties.Property
	Add(frame *model.Frame, ctx *context.TomContext) error
	AddSum(sum PropertyValueSum, ctx *context.TomContext) error
	FormatRounded(ctx *context.TomContext) []string
	FormatExact(ctx *context.TomContext) []string
}

type currencyPropertySum struct {
	property *properties.Property
	rounding util.RoundingConfig

	// maps currencyCode -> map[amount]duration, we multiply at the end to reduce rounding errors
	exact   map[string]map[int64]time.Duration
	rounded map[string]map[int64]time.Duration
}

func (c *currencyPropertySum) Property() *properties.Property {
	return c.property
}

func (c *currencyPropertySum) addValue(target map[string]map[int64]time.Duration, prop *money.Money, value time.Duration) {
	if _, ok := target[prop.Currency().Code]; !ok {
		target[prop.Currency().Code] = make(map[int64]time.Duration)
	}
	target[prop.Currency().Code][prop.Amount()] += value
}

func (c *currencyPropertySum) addValues(target map[string]map[int64]time.Duration, source map[string]map[int64]time.Duration) {
	for code, values := range source {
		for amount, duration := range values {
			target[code][amount] += duration
		}
	}
}

func (c *currencyPropertySum) Add(frame *model.Frame, ctx *context.TomContext) error {
	propValue, err := ctx.Query.FindPropertyValue(c.property.ID, frame.ProjectId)
	if err == nil {
		cValue, err := propValue.AsCurrency()
		if err != nil {
			return err
		}

		c.addValue(c.exact, cValue, frame.Duration())
		c.addValue(c.rounded, cValue, util.RoundDuration(frame.Duration(), c.rounding))
	}
	return nil
}

func (c *currencyPropertySum) AddSum(sum PropertyValueSum, ctx *context.TomContext) error {
	if cSum, ok := sum.(*currencyPropertySum); !ok {
		return fmt.Errorf("property type not matching")
	} else {
		c.addValues(c.rounded, cSum.rounded)
		c.addValues(c.exact, cSum.exact)
		return nil
	}
}

func (c *currencyPropertySum) FormatRounded(ctx *context.TomContext) []string {
	return c.format(c.rounded, ctx)
}

func (c *currencyPropertySum) FormatExact(ctx *context.TomContext) []string {
	return c.format(c.exact, ctx)
}

func (c *currencyPropertySum) format(data map[string]map[int64]time.Duration, ctx *context.TomContext) []string {
	var result []string

	// var err error
	for code, values := range data {
		m := money.New(0, code)

		for amount, duration := range values {
			m, _ = m.Add(money.New(amount, code).Multiply(int64(duration.Hours())))
		}

		result = append(result, m.Display())
	}

	return result
}
