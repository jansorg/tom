package report

import (
	"time"

	"github.com/jansorg/tom/go-tom/util"
)

type Config struct {
	Splitting          []SplitOperation    `json:"splitting"`
	Timezone           *time.Location      `json:"timezone,omitempty"`
	ProjectIDs         []string            `json:"projects"`
	IncludeSubprojects bool                `json:"show_subprojects"`
	DateFilterRange    util.DateRange      `json:"date_range"`
	ShowEmpty          bool                `json:"show_empty"`
	EntryRounding      util.RoundingConfig `json:"entry_rounding,omitempty"`
	SumRounding        util.RoundingConfig `json:"sum_rounding,omitempty"`
}
