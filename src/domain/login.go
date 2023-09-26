package domain

type Login struct {
	ID       string `json:"_id,omitempty" example:"training/123"`
	Provider string `json:"provider,omitempty" example:"facebook"`
	Subject  string `json:"subject,omitempty" example:"10150000000001234"`
}

func (t Login) GetID() string {
	return t.ID
}

func (t *Login) SetID(id string) {
	t.ID = id
}
