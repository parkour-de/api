package domain

// TrainingDTO
// @description Enriches Training with some other related nodes such as Location or User
type TrainingDTO struct {
	Training
	Location     *Location `json:"location,omitempty"`
	LocationId   string    `json:"locationId,omitempty" example:"location/123"`
	OrganiserIds []string  `json:"organiserIds,omitempty" example:"user/123"`
	Organisers   []User    `json:"organisers,omitempty"`
}
