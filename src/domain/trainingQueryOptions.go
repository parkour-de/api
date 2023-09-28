package domain

// TrainingQueryOptions carries query options filtering the list of trainings or limiting the returned items or details
type TrainingQueryOptions struct {
	City        string
	Weekday     int
	OrganiserID string
	LocationID  string
	Include     map[string]struct{}
	Skip        int
	Limit       int
}
