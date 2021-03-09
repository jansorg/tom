package report

import (
	"encoding/json"
	"fmt"
	"github.com/jansorg/tom/go-tom/dateTime"
	"time"
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
	SumRounding        dateTime.RoundingConfig `json:"rounding_total"`
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
	loc, err := time.LoadLocation(string(*t))
	if err != nil {
		panic(fmt.Sprintf("unable to load timezone with name %s", string(*t)))
	}
	return loc
}
