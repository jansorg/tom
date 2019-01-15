package model

import (
	"sort"
	"strings"
)

type ProjectList []*Project

func (p ProjectList) Empty() bool {
	return len(p) == 0
}

func (p ProjectList) Size() int {
	return len(p)
}

func (p ProjectList) Projects() []*Project {
	return p
}

func (p ProjectList) First() *Project {
	if len(p) == 0 {
		return nil
	}
	return p[0]
}

func (p ProjectList) Last() *Project {
	if p.Empty() {
		return nil
	}
	return p[len(p)-1]
}

func (p ProjectList) SortByFullname() {
	sort.Slice(p, func(i, j int) bool {
		return strings.Compare(p[i].GetFullName(""), p[j].GetFullName("")) < 0
	})
}
