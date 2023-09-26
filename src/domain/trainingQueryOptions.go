package domain

// TrainingQueryOptions
// @description Carries query options filtering the list of trainings or limiting the returned items or details
type TrainingQueryOptions struct {
	City        string
	Weekday     int
	OrganiserID string
	LocationID  string
	Include     map[string]struct{}
	Skip        int
	Limit       int
}
