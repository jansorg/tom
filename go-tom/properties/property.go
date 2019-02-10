package properties

type Property struct {
	ID                 string `json:"id,required"`
	Name               string `json:"name,required"`
	TypeID             string `json:"type,required"`
	Description        string `json:"description,omitempty"`
	ApplyToSubprojects bool   `json:"subprojects,omitempty"`
}

func (p *Property) Type() PropertyType {
	return FindType(p.TypeID)
}
