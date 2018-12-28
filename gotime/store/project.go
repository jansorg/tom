package store

type Project struct {
	ID       string `json:"id"`
	ParentID string `json:"parent"`
	Name     string `json:"name"`

	FullName string `json:"fullname"`
}
