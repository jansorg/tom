package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jansorg/tom/go-tom/properties"
)

var errPropValueNotFound = fmt.Errorf("property value not found")

type Project struct {
	Store Store `json:"-"`

	ID         string                    `json:"id"`
	ParentID   string                    `json:"parent"`
	Name       string                    `json:"name"`
	Properties properties.PropertyValues `json:"properties,omitempty"`

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

func (p *Project) GetPropertyValue(id string) (properties.PropertyValue, error) {
	for _, p := range p.Properties {
		if p.PropertyID() == id {
			return p, nil
		}
	}
	return nil, errPropValueNotFound
}

func (p *Project) HasPropertyValue(id string) bool {
	_, err := p.GetPropertyValue(id)
	return err == nil
}

func (p *Project) SetPropertyValue(id string, value string) error {
	prop, err := p.Store.GetProperty(id)
	if err != nil {
		return err
	}

	propValue, err := prop.Type().Parse(value, id)
	if err != nil {
		return err
	}

	p.Properties = append(p.Properties, propValue)
	return nil
}

func (p *Project) SetProperty(value properties.PropertyValue) error {
	_, err := p.Store.GetProperty(value.PropertyID())
	if err != nil {
		return err
	}

	p.Properties = append(p.Properties, value)
	return nil
}

func (p *Project) RemovePropertyValue(id string) {
	for i, prop := range p.Properties {
		if prop.PropertyID() == id {
			p.Properties = append(p.Properties[:i], p.Properties[i+1:]...)
			return
		}
	}
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
