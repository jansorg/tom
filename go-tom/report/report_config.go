package report

import (
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
