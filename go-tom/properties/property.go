package properties

type Property struct {
	ID                 string `json:"id,required"`
	Name               string `json:"name,required"`
	Description        string `json:"description"`
	TypeID             string `json:"type"`
	ApplyToSubprojects bool   `json:"subprojects"`
}

func (p *Property) Type() PropertyType {
	return FindType(p.TypeID)
}
