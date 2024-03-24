package domain

// Training stores information about a training
type Training struct {
	Entity
	Type         string            `json:"type,omitempty" example:"training"` // parkour-training, parkour-jam, meeting, show, competition, slackline, tour
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos
	Comments []Comment `json:"comments,omitempty"`
	Cycles   []Cycle   `json:"cycles,omitempty"`
}
