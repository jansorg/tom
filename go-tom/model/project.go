package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

var errPropNotFound = fmt.Errorf("property not found")
var errPropValueNotFound = fmt.Errorf("property value not found")

type Project struct {
	Store Store `json:"-"`

	ID         string                 `json:"id"`
	ParentID   string                 `json:"parent"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties,omitempty"`

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

func (p *Project) GetPropertyValue(id string) (interface{}, error) {
	v, ok := p.Properties[id]
	if !ok {
		return nil, errPropValueNotFound
	}
	return v, nil
}

func (p *Project) HasPropertyValue(id string) bool {
	_, exists := p.Properties[id]
	return exists
}

func (p *Project) SetPropertyValue(id string, value interface{}) error {
	prop, err := p.Store.GetProperty(id)
	if err != nil {
		return err
	}

	if err := prop.Validate(value); err != nil {
		return err
	}

	p.Properties[id] = value
	return nil
}

func (p *Project) RemovePropertyValue(id string) {
	_ = p.SetPropertyValue(id, nil)
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
