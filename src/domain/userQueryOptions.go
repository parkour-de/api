package domain

// UserQueryOptions carries query options filtering the list of users or limiting the returned items or details
type UserQueryOptions struct {
	Key      string
	Name     string
	Type     string
	Text     string
	Language string
	Include  map[string]struct{}
	Skip     int
	Limit    int
}
