package util

import (
	"encoding/json"
	"fmt"
	"time"
)

var errUnknownRounding = fmt.Errorf("unknown rounding mode")

type RoundingMode int8

const (
	RoundNone RoundingMode = iota + 1
	RoundNearest
	RoundUp
)

func (r RoundingMode) String() string {
	switch r {
	case RoundNone:
		return "none"
	case RoundNearest:
		return "nearest"
	case RoundUp:
		return "up"
	default:
		return ""
	}
}

func (r RoundingMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *RoundingMode) UnmarshalJSON(data []byte) error {
	name := ""
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}

	*r = RoundingByName(name)
	return nil
}

func RoundingByName(name string) RoundingMode {
	switch name {
	case "nearest":
		return RoundNearest
	case "up":
		return RoundUp
	default:
		return RoundNone
	}
}

func RoundingNone() RoundingConfig {
	return RoundingConfig{
		Mode: RoundUp,
		Size: 0,
	}
}

func RoundingUp(size time.Duration) RoundingConfig {
	return RoundingConfig{
		Mode: RoundUp,
		Size: size,
	}
}

func RoundingNearest(size time.Duration) RoundingConfig {
	return RoundingConfig{
		Mode: RoundNearest,
		Size: size,
	}
}

type RoundingConfig struct {
	Mode RoundingMode  `json:"mode"`
	Size time.Duration `json:"size"`
}

func (r RoundingConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Mode string `json:"mode"`
		Size string `json:"size"`
	}{
		Mode: r.Mode.String(),
		Size: r.Size.String(),
	})
}

func (r *RoundingConfig) UnmarshalJSON(data []byte) error {
	d := &struct {
		Mode string `json:"mode"`
		Size string `json:"size"`
	}{}

	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}

	size, err := time.ParseDuration(d.Size)
	if err != nil {
		return err
	}

	r.Mode = RoundingByName(d.Mode)
	r.Size = size
	return nil
}
