package domain

// LocationQueryOptions carries query options filtering the list of locations or limiting the returned items or details
type LocationQueryOptions struct {
	Lat         float64 // Latitude
	Lng         float64 // Longitude
	MaxDistance float64 // Maximum distance in meters
	Type        string  // Location type filter
	Text        string
	Language    string
	Include     map[string]struct{}
	Skip        int // Skip a number of results
	Limit       int // Limit the number of results
}
