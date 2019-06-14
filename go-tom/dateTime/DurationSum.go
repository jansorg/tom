package dateTime

import (
	"time"
)

func NewDurationSum() *DurationSum {
	return &DurationSum{}
}

func NewDurationSumWithRef(referenceTime *time.Time) *DurationSum {
	return &DurationSum{
		referenceTime: referenceTime,
	}
}

func NewEmptyCopy(proto *DurationSum) *DurationSum {
	return NewDurationSumAll(proto.rounding, proto.acceptedRange, proto.referenceTime)
}

func NewDurationSumFiltered(acceptedRange *DateRange, referenceTime *time.Time) *DurationSum {
	return &DurationSum{
		acceptedRange: acceptedRange,
		referenceTime: referenceTime,
	}
}

func NewDurationSumAll(rounding RoundingConfig, acceptedRange *DateRange, referenceTime *time.Time) *DurationSum {
	return &DurationSum{
		rounding:      rounding,
		acceptedRange: acceptedRange,
		referenceTime: referenceTime,
	}
}

type DurationSum struct {
	referenceTime *time.Time
	SumRounded    time.Duration `json:"sum_rounded"`
	SumExact      time.Duration `json:"sum_exact"`
	rounding      RoundingConfig
	roundingSize  time.Duration
	acceptedRange *DateRange
}

func (d *DurationSum) IsZero() bool {
	return d.GetExact() == 0
}

func (d *DurationSum) IsRoundedZero() bool {
	return d.Get() == 0
}

func (d *DurationSum) IsRounded() bool {
	return d.SumExact != d.SumRounded
}

func (d *DurationSum) AddSum(r *DurationSum) {
	// fixme handle incompatible config values of r?
	d.SumExact += r.SumExact
	d.SumRounded += r.SumRounded
}

func (d *DurationSum) AddRange(r DateRange) {
	d.add(r.Start, r.End)
}

func (d *DurationSum) AddStartEnd(start time.Time, end time.Time) {
	d.add(&start, &end)
}

func (d *DurationSum) AddStartEndP(start *time.Time, end *time.Time) {
	if end == nil {
		d.add(start, d.referenceTime)
	} else {
		d.add(start, end)
	}
}

func (d *DurationSum) Add(duration time.Duration) {
	d.SumExact += duration
	d.SumRounded += RoundDuration(duration, d.rounding)
}

func (d *DurationSum) Get() time.Duration {
	return d.SumRounded
}

func (d *DurationSum) GetExact() time.Duration {
	return d.SumExact
}

func (d *DurationSum) add(a *time.Time, b *time.Time) {
	if b == nil && d.referenceTime == nil {
		return
	}

	start := a
	end := b
	if end == nil {
		end = d.referenceTime
	}

	var overlap time.Duration
	if d.acceptedRange != nil {
		overlap = d.acceptedRange.Intersection(start, end)
	} else {
		overlap = b.Sub(*a)
	}

	d.Add(overlap)
}
