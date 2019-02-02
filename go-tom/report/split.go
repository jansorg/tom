package report

import (
	"encoding/json"
	"fmt"
)

type SplitOperation int8

const (
	SplitByYear SplitOperation = iota + 1
	SplitByMonth
	SplitByWeek
	SplitByDay
	SplitByProject
)

func SplitOperationByName(name string) (SplitOperation, error) {
	switch name {
	case "year":
		return SplitByYear, nil
	case "month":
		return SplitByMonth, nil
	case "week":
		return SplitByWeek, nil
	case "day":
		return SplitByDay, nil
	case "project":
		return SplitByProject, nil
	default:
		return 0, fmt.Errorf("unknown split operation %s", name)
	}
}

func (s SplitOperation) IsDateSplit() bool {
	return s >= SplitByYear && s <= SplitByDay
}

func (s SplitOperation) String() string {
	name := ""
	switch (s) {
	case SplitByYear:
		name = "year"
	case SplitByMonth:
		name = "month"
	case SplitByWeek:
		name = "week"
	case SplitByDay:
		name = "day"
	case SplitByProject:
		name = "project"
	}
	return name
}
func (s SplitOperation) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (r *SplitOperation) UnmarshalJSON(data []byte) error {
	name := ""
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}

	v, err := SplitOperationByName(name)
	if err != nil {
		return err
	}
	*r = v
	return nil
}
