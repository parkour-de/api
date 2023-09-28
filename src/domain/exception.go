package domain

import "time"

// Exception to a recurring event
type Exception struct {
	Date       time.Time `json:"date,omitempty"`     // RFC 3339 date of when the exception occurs
	Begin      int       `json:"begin,omitempty"`    // seconds
	Duration   int       `json:"duration,omitempty"` // seconds, use 0 to cancel
	LocationId string    `json:"locationId,omitempty" example:"location/123"`
}
