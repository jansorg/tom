package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jansorg/tom/go-tom/money"
)

var errPropValueNotFound = fmt.Errorf("property value not found")

type ProjectProperties struct {
	HourlyRate *money.Money `json:"hourlyRate,omitempty"`
}

type Project struct {
	ID         string             `json:"id"`
	ParentID   string             `json:"parent"`
	Name       string             `json:"name"`
	Properties *ProjectProperties `json:"properties,omitempty"`

	Store    Store    `json:"-"`
	FullName []string `json:"-"`
}

func (p *Project) GetFullName(delimiter string) string {
	return strings.Join(p.FullName, delimiter)
}

func (p *Project) Parent() *Project {
	if p.ParentID == "" {
		return nil
	}

	parent, _ := p.Store.ProjectByID(p.ParentID)
	return parent
}

func (p *Project) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("project id must not be empty")
	}
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("project name must not be empty")
	}
	return nil
}

func (p *Project) HourlyRate() *money.Money {
	if p.Properties == nil {
		return nil
	}
	return p.Properties.HourlyRate
}

func (p *Project) SetHourlyRate(value *money.Money) {
	if value == nil {
		if p.Properties != nil {
			// fixme only valid as long as we have a single property
			p.Properties = nil
		}
		return
	}

	if p.Properties == nil {
		p.Properties = &ProjectProperties{}
	}
	p.Properties.HourlyRate = value
}

type DetailedProject Project

func (p *DetailedProject) MarshalJSON() ([]byte, error) {
	type Alias DetailedProject
	return json.Marshal(&struct {
		FullName []string `json:"fullName"`
		*Alias
	}{
		FullName: p.FullName,
		Alias:    (*Alias)(p),
	})
}
