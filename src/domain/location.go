package domain

type Location struct {
	Entity
	Lat          float64           `json:"lat,omitempty" example:"53.55"`
	Lng          float64           `json:"lng,omitempty" example:"9.99"`
	City         string            `json:"city,omitempty" example:"Hamburg"`
	Type         string            `json:"type,omitempty" example:"spot"` // spot, gym, parkour-gym, office, public-transport
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos
	Comments []Comment `json:"comments,omitempty"`
}
