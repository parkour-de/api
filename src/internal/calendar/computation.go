package calendar

import (
	"pkv/api/src/domain"
	"time"
)

// Occurrence represents an instance of a cycle.
type Occurrence struct {
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

// ComputeOccurrences computes occurrences within the specified date range.
func ComputeOccurrences(cycle domain.Cycle, start, end time.Time) []time.Time {
	if cycle.Interval == 0 {
		cycle.Interval = 1
	}

	var occurrences []time.Time

	// Initialize the current date to the first satisfying date on or after the start date.
	current := fixWeekday(cycle, cycle.Startdate)
	for current.Before(start) {
		current = ComputeNextSatisfyingDate(cycle, current)
	}

	// Compute occurrences within the specified date range.
	for current.Before(end) && current.Before(cycle.Enddate) {
		occurrences = append(occurrences, current)
		current = ComputeNextSatisfyingDate(cycle, current)
	}

	return occurrences
}

func fixWeekday(cycle domain.Cycle, current time.Time) time.Time {
	if cycle.Monthday == 0 && cycle.Weekday != 0 {
		daysToNextWeekday := (cycle.Weekday - int(current.Weekday()) + 7) % 7
		current = current.AddDate(0, 0, daysToNextWeekday)
	} else if cycle.Monthday != 0 {
		dayInMonth := calculateDayInMonth(current.Year(), current.Month(), cycle.Monthday, cycle.Weekday)
		if dayInMonth < current.Day() {
			dayInMonth = calculateDayInMonth(current.Year(), current.Month()+1, cycle.Monthday, cycle.Weekday)
			current = time.Date(current.Year(), current.Month()+1, dayInMonth, 0, 0, 0, 0, current.Location())
		} else {
			current = time.Date(current.Year(), current.Month(), dayInMonth, 0, 0, 0, 0, current.Location())
		}
	}
	return current
}

func ComputeNextSatisfyingDate(cycle domain.Cycle, current time.Time) time.Time {
	if cycle.Monthday == 0 && cycle.Weekday == 0 {
		// Advance by Interval days.
		current = current.AddDate(0, 0, cycle.Interval)
	} else if cycle.Monthday == 0 {
		// Advance by Interval weeks.
		current = current.AddDate(0, 0, 7*cycle.Interval)
	} else {
		// Advance by Interval months.
		nextMonth := time.Date(current.Year(), current.Month()+time.Month(cycle.Interval), 1, 0, 0, 0, 0, current.Location())

		// Calculate the day within the month based on Monthday and Weekday.
		dayInMonth := calculateDayInMonth(nextMonth.Year(), nextMonth.Month(), cycle.Monthday, cycle.Weekday)

		// Update the current date.
		current = time.Date(nextMonth.Year(), nextMonth.Month(), dayInMonth, 0, 0, 0, 0, current.Location())
	}

	return current
}

func calculateDayInMonth(year int, month time.Month, monthday, weekday int) int {
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	normalisedWeekday := weekday
	if normalisedWeekday == 7 {
		normalisedWeekday = 0
	}

	if monthday > 0 {
		// Count days from the beginning of the month.
		dayInMonth := monthday
		if weekday != 0 {
			weekdayCount := monthday
			for i := 1; i <= daysInMonth; i++ {
				currentDay := firstOfMonth.AddDate(0, 0, i-1).Weekday()
				if int(currentDay) == normalisedWeekday {
					weekdayCount--
					if weekdayCount == 0 {
						dayInMonth = i
						break
					}
				}
			}
		}
		return dayInMonth
	} else if monthday < 0 {
		lastWeekday := 0
		if weekday != 0 {
			weekdayCount := -monthday
			for i := daysInMonth; i >= 1; i-- {
				currentDay := firstOfMonth.AddDate(0, 0, i-1).Weekday()
				if int(currentDay) == normalisedWeekday {
					weekdayCount--
					if weekdayCount == 0 {
						lastWeekday = i
						break
					}
				}
			}
		} else {
			// Count days from the end of the month.
			lastWeekday = daysInMonth + 1 + monthday
		}
		return lastWeekday
	}

	// Monthday is zero; this is unexpected. Return 0 as an invalid result
	return 0
}
