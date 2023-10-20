package domain

import "time"

type Location struct {
	Key          string            `json:"_key,omitempty" example:"123"`
	Created      time.Time         `json:"created,omitempty"`  // RFC 3339 date
	Modified     time.Time         `json:"modified,omitempty"` // RFC 3339 date
	Lat          float64           `json:"lat,omitempty" example:"53.55"`
	Lng          float64           `json:"lng,omitempty" example:"9.99"`
	City         string            `json:"city,omitempty" example:"Hamburg"`
	Type         string            `json:"type,omitempty" example:"spot"` // spot, gym, parkour-gym, office, public-transport
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos       []Photo           `json:"photos,omitempty"`
	Comments     []Comment         `json:"comments,omitempty"`
}

func (l Location) GetKey() string {
	return l.Key
}

func (l *Location) SetKey(id string) {
	l.Key = id
}
