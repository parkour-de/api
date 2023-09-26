package domain

type Location struct {
	ID           string            `json:"_id,omitempty" example:"location/123"`
	Lat          float64           `json:"lat,omitempty" example:"53.55"`
	Lng          float64           `json:"lng,omitempty" example:"9.99"`
	City         string            `json:"city,omitempty" example:"Hamburg"`
	Type         string            `json:"type,omitempty" example:"spot"` // spot, gym, parkour-gym, office, public-transport
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos       []Photo           `json:"photos,omitempty"`
	Comments     []Comment         `json:"comments,omitempty"`
}

func (l Location) GetID() string {
	return l.ID
}

func (l *Location) SetID(id string) {
	l.ID = id
}
