package model

import (
	"fmt"
	"strings"
)

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (t *Tag) Validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return fmt.Errorf("tag name must not be empty")
	}
	return nil
}
