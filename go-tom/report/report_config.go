package report

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jansorg/tom/go-tom/dateTime"
)

type Config struct {
	Splitting          []SplitOperation        `json:"split"`
	TimezoneName       TimezoneName            `json:"timezone"`
	ProjectIDs         []string                `json:"projects"`
	ProjectDelimiter   string                  `json:"project_delimiter"`
	IncludeSubprojects bool                    `json:"show_subprojects"`
	DateFilterRange    dateTime.DateRange      `json:"date_range"`
	ShowEmpty          bool                    `json:"show_empty"`
	ShowStopTime       bool                    `json:"show_stop_time"`
	IncludeArchived    bool                    `json:"include_archived"`
	ShortTitles        bool                    `json:"short_titles"`
	EntryRounding      dateTime.RoundingConfig `json:"rounding_entry"`
}

func NewTimezoneName(location *time.Location) TimezoneName {
	return TimezoneName(location.String())
}

func NewTimezoneNameUTC() TimezoneName {
	return NewTimezoneName(time.UTC)
}

func NewTimezoneNameLocal() TimezoneName {
	return NewTimezoneName(time.Local)
}

type TimezoneName string

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
