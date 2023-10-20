package domain

import "time"

// Page stores information about a page
type Page struct {
	Key          string            `json:"_key,omitempty" example:"123"`
	Created      time.Time         `json:"created,omitempty"`  // RFC 3339 date
	Modified     time.Time         `json:"modified,omitempty"` // RFC 3339 date
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos       []Photo           `json:"photos,omitempty"`
	Comments     []Comment         `json:"comments,omitempty"`
}

func (p Page) GetKey() string {
	return p.Key
}

func (p *Page) SetKey(id string) {
	p.Key = id
}
