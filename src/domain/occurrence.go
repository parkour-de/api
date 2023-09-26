package domain

import "time"

// Occurrence
// @description Automatically generated individual calendar entries based on Cycle and Exception
type Occurrence struct {
	Date       time.Time `json:"date,omitempty"`     // RFC 3339 date of when it occurs
	Begin      int       `json:"begin,omitempty"`    // seconds
	Duration   int       `json:"duration,omitempty"` // seconds
	LocationId string    `json:"locationId,omitempty" example:"location/123"`
}
