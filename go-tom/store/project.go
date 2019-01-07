package store

type Project struct {
	store Store

	ID       string `json:"id"`
	ParentID string `json:"parent"`
	Name     string `json:"name"`

	FullName string `json:"fullname"`

	Properties map[string]string `json:"properties,omitempty"`
}

func (p *Project) GetProperties() map[string]string {
	return p.Properties
}

// func (p *Project) setProperty(key string, value interface{}) error {
// 	p.Properties[key] = value
// 	if newP, err := p.store.UpdateProject(p); err != nil {
// 		return err
// 	} else {
// 		*p = *newP
// 		return nil
// 	}
// }
//
// func (p *Project) getString(key string) (string, bool) {
// 	if p, ok := p.Properties[key]; !ok {
// 		return "", false
// 	} else {
// 		s, ok := p.(string)
// 		return s, ok
// 	}
// }
//
// func (p *Project) getBool(key string) (bool, bool) {
// 	if p, ok := p.Properties[key]; !ok {
// 		return false, false
// 	} else {
// 		s, ok := p.(bool)
// 		return s, ok
// 	}
// }
//
// func (p *Project) getInt64(key string) (int64, bool) {
// 	if p, ok := p.Properties[key]; !ok {
// 		return 0, false
// 	} else {
// 		s, ok := p.(int64)
// 		return s, ok
// 	}
// }
//
// func (p *Project) getFloat64(key string) (float64, bool) {
// 	if p, ok := p.Properties[key]; !ok {
// 		return 0, false
// 	} else {
// 		s, ok := p.(float64)
// 		return s, ok
// 	}
// }
