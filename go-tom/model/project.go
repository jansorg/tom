package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Project struct {
	Store Store `json:"-"`

	ID       string `json:"id"`
	ParentID string `json:"parent"`
	Name     string `json:"name"`

	FullName []string `json:"-"`

	Properties map[string]string `json:"properties,omitempty"`
}

func (p *Project) GetFullName(delimiter string) string {
	return strings.Join(p.FullName, delimiter)
}

func (p *Project) GetProperties() map[string]string {
	return p.Properties
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
