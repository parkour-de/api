package domain

// TrainingDTO enriches Training with some other related nodes such as Location or User
type TrainingDTO struct {
	Training
	Location      *Location `json:"location,omitempty"`
	LocationKey   string    `json:"locationId,omitempty" example:"123"`
	OrganiserKeys []string  `json:"organiserIds,omitempty" example:"123"`
	Organisers    []User    `json:"organisers,omitempty"`
}
