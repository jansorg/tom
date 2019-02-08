package report

import (
	"time"

	"github.com/jansorg/tom/go-tom/util"
)

type Config struct {
	Splitting          []SplitOperation    `json:"split"`
	Timezone           *time.Location      `json:"timezone"`
	ProjectIDs         []string            `json:"projects"`
	IncludeSubprojects bool                `json:"show_subprojects"`
	DateFilterRange    util.DateRange      `json:"date_range"`
	ShowEmpty          bool                `json:"show_empty"`
	IncludeArchived    bool                `json:"include_archived"`
	EntryRounding      util.RoundingConfig `json:"rounding_entry"`
	SumRounding        util.RoundingConfig `json:"rounding_total"`
}
