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

	if _, err := parseTimezone(name); err != nil {
		return err
	}

	*t = TimezoneName(name)
	return nil
}

func (t *TimezoneName) AsTimezone() *time.Location {
	loc, err := parseTimezone(string(*t))
	if err != nil {
		panic(err.Error())
	}
	return loc
}

func parseTimezone(name string) (*time.Location, error) {
	var offset int
	if n, err := fmt.Sscanf(name, "UTC+%d", &offset); err == nil && n == 1 {
		return time.FixedZone(name, offset*60*60), nil
	}

	if n, err := fmt.Sscanf(name, "UTC-%d", &offset); err == nil && n == 1 {
		return time.FixedZone(name, -offset*60*60), nil
	}

	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("unable to load timezone with name %s: %v", name, err)
	}
	return loc, nil
}
