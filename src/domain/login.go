package domain

import "time"

type Login struct {
	Key      string    `json:"_key,omitempty" example:"123"`
	Provider string    `json:"provider,omitempty" example:"facebook"`
	Subject  string    `json:"subject,omitempty" example:"10150000000001234"`
	Enabled  bool      `json:"enabled,omitempty"`
	Created  time.Time `json:"created,omitempty"` // RFC 3339 date
}

func (t Login) GetKey() string {
	return t.Key
}

func (t *Login) SetKey(id string) {
	t.Key = id
}
