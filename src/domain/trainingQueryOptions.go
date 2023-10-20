package domain

// TrainingQueryOptions carries query options filtering the list of trainings or limiting the returned items or details
type TrainingQueryOptions struct {
	City         string
	Weekday      int
	OrganiserKey string
	LocationKey  string
	Text         string
	Language     string
	Include      map[string]struct{}
	Skip         int
	Limit        int
}
