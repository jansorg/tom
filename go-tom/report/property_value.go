package report

import (
	"github.com/jansorg/tom/go-tom/properties"
)

type PropertyValue struct {
	Property        *properties.Property `json:"property"`
	ValueForRounded float64              `json:"value_rounded"`
	ValueForExact   float64              `json:"value_exact"`
}
