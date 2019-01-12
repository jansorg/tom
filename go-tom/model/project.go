package model

import "strings"

type Project struct {
	Store Store `json:"-"`

	ID       string   `json:"id"`
	ParentID string   `json:"parent"`
	Name     string   `json:"name"`
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