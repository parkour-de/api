package domain

import "time"

type Login struct {
	Key      string    `json:"_key,omitempty" example:"123"`
	Created  time.Time `json:"created,omitempty"`  // RFC 3339 date
	Modified time.Time `json:"modified,omitempty"` // RFC 3339 date
	Provider string    `json:"provider,omitempty" example:"facebook"`
	Subject  string    `json:"subject,omitempty" example:"10150000000001234"`
	Enabled  bool      `json:"enabled,omitempty"`
}

func (t Login) GetKey() string {
	return t.Key
}

func (t *Login) SetKey(id string) {
	t.Key = id
}
