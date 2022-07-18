package report

import (
	"encoding/json"
	"fmt"
	"time"
)

type TimezoneName string

func NewTimezoneName(location *time.Location) TimezoneName {
	return TimezoneName(location.String())
}

func NewTimezoneNameUTC() TimezoneName {
	return NewTimezoneName(time.UTC)
}

func NewTimezoneNameLocal() TimezoneName {
	return NewTimezoneName(time.Local)
}

func (t *TimezoneName) UnmarshalJSON(bytes []byte) error {
	var name string
	if err := json.Unmarshal(bytes, &name); err != nil {
		return err
	}

	if _, err := time.LoadLocation(name); err != nil {
		panic(fmt.Sprintf("error locating time zone %s: %v", name, err))
		return err
	}

	*t = TimezoneName(name)
	return nil
}

func (t *TimezoneName) AsTimezone() *time.Location {
	name := string(*t)

	var offset int
	if n, err := fmt.Sscanf(name, "UTC+%d", &offset); err == nil && n == 1 {
		return time.FixedZone(name, offset*60*60)
	}

	if n, err := fmt.Sscanf(name, "UTC-%d", &offset); err == nil && n == 1 {
		return time.FixedZone(name, -offset*60*60)
	}

	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("unable to load timezone with name %s", name))
	}
	return loc
}
