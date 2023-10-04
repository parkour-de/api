package domain

// LocationDTO is a location with distance to a given point
type LocationDTO struct {
	Location
	Distance float64 `json:"distance,omitempty"`
}
