package calendar

import (
	"pkv/api/src/domain"
	"sort"
	"time"
)

// ComputeDays computes occurrences within the specified date range.
func ComputeDays(cycle domain.Cycle, start, end time.Time) []time.Time {
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

// GenerateOccurrences is a convenience function for ComputeDays. ComputeDays did
// compute a list of time.Time values with the same location as the start
// parameter and the hour, min, sec set to zero. However, cycle contains the
// beginning and duration as seconds from midnight (and we assume the current
// timezone is correctly chosen). The following function will convert the return
// value of ComputeDays into a list of domain.Occurrence:
func GenerateOccurrences(cycle domain.Cycle, occurrences []time.Time) []domain.Occurrence {
	var newOccurrences []domain.Occurrence
	for _, occurrence := range occurrences {
		newOccurrences = append(newOccurrences, domain.Occurrence{
			Date:       occurrence,
			Begin:      cycle.Begin,
			Duration:   cycle.Duration,
			LocationId: cycle.LocationId,
		})
	}
	return newOccurrences
}

// ApplyExceptions is a convenience function for GenerateOccurrences. Now, a
// Cycle with Occurrences can also have Exceptions. So after generating some
// Occurrences for the cycle and provided timespan, the list of Exception structs
// given will provide either extra events (on days where ComputeDays hadn't had
// an event), events with changed location or time (when the duration is greater
// than zero) or a cancelled event (when the duration is zero). The following
// function taking the list of Occurrences and Exceptions and a start and end day
// again will return a list of Occurrences:
func ApplyExceptions(occurrences []domain.Occurrence, exceptions []domain.Exception, start, end time.Time) []domain.Occurrence {
	var newOccurrences []domain.Occurrence
nextOccurrence:
	for _, occurrence := range occurrences {
		// Skip if an exception exists for the current occurrence.
		for _, exception := range exceptions {
			if exception.Date.Equal(occurrence.Date) {
				continue nextOccurrence
			}
		}
		newOccurrences = append(newOccurrences, occurrence)
	}
	// Add the exceptions that are not already in the list of occurrences.
	for _, exception := range exceptions {
		if exception.Date.Before(start) || exception.Date.After(end) {
			continue
		}
		// Check if the exception is a cancellation.
		if exception.Duration == 0 {
			continue
		}
		// Add the exception to the list of occurrences.
		newOccurrences = append(newOccurrences, domain.Occurrence{
			Date:       exception.Date,
			Begin:      exception.Begin,
			Duration:   exception.Duration,
			LocationId: exception.LocationId,
		})
	}
	// finally we need to sort the occurrences by date and by begin:
	sort.Slice(newOccurrences, func(i, j int) bool {
		if newOccurrences[i].Date.Equal(newOccurrences[j].Date) {
			return newOccurrences[i].Begin < newOccurrences[j].Begin
		}
		return newOccurrences[i].Date.Before(newOccurrences[j].Date)
	})
	return newOccurrences
}

// TrimOccurrences helps to clean up the list of occurrences stored in a database.

func trimOccurrences(occurrences []domain.Occurrence, start, end time.Time) []domain.Occurrence {
	var newOccurrences []domain.Occurrence
	for _, occurrence := range occurrences {
		if occurrence.Date.Before(start) || occurrence.Date.After(end) {
			continue
		}
		newOccurrences = append(newOccurrences, occurrence)
	}
	return newOccurrences
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
