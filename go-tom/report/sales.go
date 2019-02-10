package report

import (
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/money"
	"github.com/jansorg/tom/go-tom/util"
)

func NewSales(ctx *context.TomContext, entryRounding util.RoundingConfig) *Sales {
	return &Sales{
		ctx:           ctx,
		entryRounding: entryRounding,
		exactValues:   make(map[string]*money.Money),
		values:        make(map[string]*money.Money),
	}
}

type Sales struct {
	ctx           *context.TomContext
	exactValues   map[string]*money.Money
	values        map[string]*money.Money
	entryRounding util.RoundingConfig
	sumRounding   util.RoundingConfig
}

func (s *Sales) Add(frame *model.Frame) error {
	hourlyRate, err := s.ctx.Query.HourlyRate(frame.ProjectId)
	if err == nil {
		duration := frame.Duration()
		rounded := util.RoundDuration(duration, s.entryRounding)

		if _, ok := s.exactValues[hourlyRate.CurrencyCode()]; !ok {
			s.exactValues[hourlyRate.CurrencyCode()] = money.NewMoney(0, hourlyRate.CurrencyCode())
		}
		if err := s.exactValues[hourlyRate.CurrencyCode()].Add(hourlyRate.Multiple(duration.Hours())); err != nil {
			return err
		}

		if _, ok := s.values[hourlyRate.CurrencyCode()]; !ok {
			s.values[hourlyRate.CurrencyCode()] = money.NewMoney(0, hourlyRate.CurrencyCode())
		}
		if err := s.values[hourlyRate.CurrencyCode()].Add(hourlyRate.Multiple(rounded.Hours())); err != nil {
			return err
		}
	}
	return nil
}

func (s *Sales) AddSales(sales *Sales) error {
	for k, v := range sales.values {
		if _, ok := s.values[k]; !ok {
			s.values[k] = money.NewMoney(0, v.CurrencyCode())
		}

		if err := s.values[k].Add(v); err != nil {
			return err
		}
	}

	for k, v := range sales.exactValues {
		if _, ok := s.exactValues[k]; !ok {
			s.exactValues[k] = money.NewMoney(0, v.CurrencyCode())
		}

		if err := s.exactValues[k].Add(v); err != nil {
			return err
		}
	}

	return nil
}

func (s *Sales) Exact() []*money.Money {
	var result []*money.Money

	for _, v := range s.exactValues {
		result = append(result, v)
	}

	return result
}

func (s *Sales) Rounded() []*money.Money {
	var result []*money.Money

	for _, v := range s.values {
		result = append(result, v)
	}

	return result
}
