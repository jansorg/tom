package model

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
	if p.Empty() {
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
