package report

import "github.com/jansorg/tom/go-tom/model"

type PropertyValue struct {
	Property        *model.Property `json:"property"`
	ValueForRounded float64         `json:"value_rounded"`
	ValueForExact   float64         `json:"value_exact"`
}
