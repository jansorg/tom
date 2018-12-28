package store

import (
	"fmt"
	"log"
	"strings"
)

type Project struct {
	store Store

	ID       string `json:"id"`
	ParentID string `json:"parent"`
	Name     string `json:"name"`
}

func (p *Project) FullName() string {
	if p.ParentID == "" {
		return p.Name
	}

	parents := []string{p.Name}

	id := p.ParentID
	for id != "" {
		parent, err := p.store.FindFirstProject(func(current *Project) bool {
			return current.ID == id
		})

		if err != nil {
			log.Fatal(fmt.Errorf("unable to find project %s", id))
		}
		id = parent.ParentID
	}

	return strings.Join(parents, "/")
}
