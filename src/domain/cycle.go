package domain

import "time"

// Cycle represents recurring events
type Cycle struct {
	Weekday    int       `json:"weekday,omitempty"`
	Monthday   int       `json:"monthday,omitempty"`
	Interval   int       `json:"interval,omitempty"`
	Begin      int       `json:"begin,omitempty"`     // seconds
	Duration   int       `json:"duration,omitempty"`  // seconds
	Startdate  time.Time `json:"startdate,omitempty"` // RFC 3339 date of the first possible day in a cycle
	Enddate    time.Time `json:"enddate,omitempty"`   // RFC 3339 date of when the cycle does not continue
	LocationId string    `json:"locationId,omitempty" example:"location/123"`
}

/*
Example values for Cycle:

Jeden Tag:
Weekday = 0
Monthday = 0
Interval = 0 | 1

Jeden Freitag:
Weekday = 5
Monthday = 0
Interval = 0 | 1

Jeden zweiten Sonntag:
Weekday = 7
Monthday = 0
Interval = 2

Jeden ersten Donnerstag im Monat:
Weekday = 4
Monthday = 1
Interval = 0 | 1

Alle zwei Monate jeden vorletzten Mittwoch:
Weekday = 3
Monthday = -2
Interval = 2

Am 3. Tag jedes Monats:
Weekday = 0
Monthday = 3
Interval = 0 | 1

Weekday = 0 kombiniert mit Interval divisible by 7 ist verboten
Trainings, die auch am Freitag stattfinden: Weekday in [0, 5]

StartDate
EndDate
Begin
Duration
LocationKey
*/
