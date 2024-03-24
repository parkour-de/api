package domain

import "time"

type Entity struct {
	Key      string    `json:"_key,omitempty" example:"123"`
	Created  time.Time `json:"created,omitempty"`  // RFC 3339 date
	Modified time.Time `json:"modified,omitempty"` // RFC 3339 date
}

func (e Entity) GetKey() string {
	return e.Key
}

func (e *Entity) SetKey(id string) {
	e.Key = id
}

func Key(key string) Entity {
	return Entity{Key: key}
}
