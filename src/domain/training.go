package domain

import "time"

// Training stores information about a training
type Training struct {
	Key          string            `json:"_key,omitempty" example:"123"`
	Created      time.Time         `json:"created,omitempty"`                 // RFC 3339 date
	Modified     time.Time         `json:"modified,omitempty"`                // RFC 3339 date
	Type         string            `json:"type,omitempty" example:"training"` // parkour-training, parkour-jam, meeting, show, competition, slackline, tour
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos       []Photo           `json:"photos,omitempty"`
	Comments     []Comment         `json:"comments,omitempty"`
	Cycles       []Cycle           `json:"cycles,omitempty"`
}

func (t Training) GetKey() string {
	return t.Key
}

func (t *Training) SetKey(id string) {
	t.Key = id
}
