package domain

type Login struct {
	Key      string `json:"_key,omitempty" example:"123"`
	Provider string `json:"provider,omitempty" example:"facebook"`
	Subject  string `json:"subject,omitempty" example:"10150000000001234"`
}

func (t Login) GetKey() string {
	return t.Key
}

func (t *Login) SetKey(id string) {
	t.Key = id
}
