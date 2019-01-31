package dateUtil

import (
	"time"
)

func NewDurationSum() *DurationSum {
	return &DurationSum{}
}

func NewDurationLike(proto *DurationSum) *DurationSum {
	return NewDurationSumAll(proto.roundingMode, proto.roundingSize, proto.acceptedRange, proto.referenceTime)
}

func NewDurationSumFiltered(acceptedRange *DateRange, referenceTime *time.Time) *DurationSum {
	return &DurationSum{
		acceptedRange: acceptedRange,
		referenceTime: referenceTime,
	}
}

func NewDurationSumAll(roundingMode RoundingMode, roundingSize time.Duration, acceptedRange *DateRange, referenceTime *time.Time) *DurationSum {
	return &DurationSum{
		roundingMode:  roundingMode,
		roundingSize:  roundingSize,
		acceptedRange: acceptedRange,
		referenceTime: referenceTime,
	}
}

type DurationSum struct {
	referenceTime *time.Time
	SumRounded    time.Duration `json:"sum_rounded"`
	SumExact      time.Duration `json:"sum_exact"`
	roundingMode  RoundingMode
	roundingSize  time.Duration
	acceptedRange *DateRange
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
	d.add(start, end)
}

func (d *DurationSum) Add(duration time.Duration) {
	d.SumExact += duration
	d.SumRounded += RoundDuration(duration, d.roundingMode, d.roundingSize)
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
