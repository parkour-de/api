package domain

// Training stores information about a training
type Training struct {
	ID           string            `json:"_id,omitempty" example:"training/123"`
	Type         string            `json:"type,omitempty" example:"training"` // parkour-training, parkour-jam, meeting, show, competition, slackline
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos       []Photo           `json:"photos,omitempty"`
	Comments     []Comment         `json:"comments,omitempty"`
	Cycles       []Cycle           `json:"cycles,omitempty"`
}

func (t Training) GetID() string {
	return t.ID
}

func (t *Training) SetID(id string) {
	t.ID = id
}
