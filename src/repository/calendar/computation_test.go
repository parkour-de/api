package calendar

import (
	"pkv/api/src/domain"
	"reflect"
	"testing"
	"time"
)

func Test_calculateDayInMonth(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    time.Month
		monthday int
		weekday  int
		want     int
	}{
		{"first of January", 2023, time.January, 1, 0, 1},
		{"third of January", 2023, time.January, 3, 0, 3},
		{"last of January", 2023, time.January, -1, 0, 31},
		{"penultimate of January", 2023, time.January, -2, 0, 30},
		{"first Monday of January", 2023, time.January, 1, 1, 2},
		{"first Sunday of January", 2023, time.January, 1, 7, 1},
		{"second Monday of January", 2023, time.January, 2, 1, 9},
		{"fifth Monday of January", 2023, time.January, 5, 1, 30},
		{"last Wednesday of January", 2023, time.January, -1, 3, 25},
		{"last Tuesday of January", 2023, time.January, -1, 2, 31},
		{"penultimate Saturday of January", 2023, time.January, -2, 6, 21},
		{"penultimate Sunday of January", 2023, time.January, -2, 7, 22},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := calculateDayInMonth(tt.year, tt.month, tt.monthday, tt.weekday)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s\ncalculateDayInMonth(%d, %d, %d, %d) = %v, want %v", tt.name, tt.year, tt.month, tt.monthday, tt.weekday, got, tt.want)
			}
		})
	}
}

func TestComputeDays(t *testing.T) {
	currentDate := time.Now()
	jan1 := time.Date(2023, 1, 1, 0, 0, 0, 0, currentDate.Location())
	dec31 := time.Date(2023, 12, 31, 0, 0, 0, 0, currentDate.Location())
	mar1 := time.Date(2023, 3, 1, 0, 0, 0, 0, currentDate.Location())
	// mar2 := time.Date(2023, 3, 2, 0, 0, 0, 0, currentDate.Location())
	mar5 := time.Date(2023, 3, 5, 0, 0, 0, 0, currentDate.Location())
	mar31 := time.Date(2023, 3, 31, 0, 0, 0, 0, currentDate.Location())
	tests := []struct {
		name  string
		cycle domain.Cycle
		start time.Time
		end   time.Time
		want  []time.Time
	}{
		{"every day",
			domain.Cycle{0, 0, 0, 1800, 900, jan1, dec31, ""},
			mar1, mar5,
			[]time.Time{
				time.Date(2023, 3, 1, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 2, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 3, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 4, 0, 0, 0, 0, currentDate.Location()),
			},
		},
		{"every day, short period",
			domain.Cycle{0, 0, 0, 1800, 900, mar1, mar5, ""},
			jan1, dec31,
			[]time.Time{
				time.Date(2023, 3, 1, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 2, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 3, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 4, 0, 0, 0, 0, currentDate.Location()),
			},
		},
		{"every Friday",
			domain.Cycle{5, 0, 0, 1800, 900, jan1, dec31, ""},
			mar1, mar31,
			[]time.Time{
				time.Date(2023, 3, 3, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 10, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 17, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 24, 0, 0, 0, 0, currentDate.Location()),
			},
		},
		{"every Friday, short period",
			domain.Cycle{5, 0, 0, 1800, 900, mar1, mar31, ""},
			jan1, dec31,
			[]time.Time{
				time.Date(2023, 3, 3, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 10, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 17, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 24, 0, 0, 0, 0, currentDate.Location()),
			},
		},
		{"every second Sunday in every second month",
			domain.Cycle{7, 2, 2, 1800, 900, jan1, dec31, ""},
			mar1, dec31,
			[]time.Time{
				time.Date(2023, 3, 12, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 5, 14, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 7, 9, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 9, 10, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 11, 12, 0, 0, 0, 0, currentDate.Location()),
			},
		},
		{"every second Sunday in every second month, short period",
			domain.Cycle{7, 2, 2, 1800, 900, mar1, dec31, ""},
			jan1, dec31,
			[]time.Time{
				time.Date(2023, 3, 12, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 5, 14, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 7, 9, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 9, 10, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 11, 12, 0, 0, 0, 0, currentDate.Location()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeDays(tt.cycle, tt.start, tt.end)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComputeDays(%#v, %v, %v)\n  got = %v,\n  want  %v", tt.cycle, tt.start, tt.end, got, tt.want)
			}
		})
	}
}

func TestGenerateOccurrences(t *testing.T) {
	currentDate := time.Now()
	jan1 := time.Date(2023, 1, 1, 0, 0, 0, 0, currentDate.Location())
	dec31 := time.Date(2023, 12, 31, 0, 0, 0, 0, currentDate.Location())
	tests := []struct {
		name  string
		cycle domain.Cycle
		given []time.Time
		want  []domain.Occurrence
	}{
		{"test",
			domain.Cycle{0, 0, 0, 1800, 900, jan1, dec31, ""},
			[]time.Time{
				time.Date(2023, 3, 1, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 2, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 3, 0, 0, 0, 0, currentDate.Location()),
				time.Date(2023, 3, 4, 0, 0, 0, 0, currentDate.Location()),
			},
			[]domain.Occurrence{
				{time.Date(2023, 3, 1, 0, 0, 0, 0, currentDate.Location()), 1800, 900, ""},
				{time.Date(2023, 3, 2, 0, 0, 0, 0, currentDate.Location()), 1800, 900, ""},
				{time.Date(2023, 3, 3, 0, 0, 0, 0, currentDate.Location()), 1800, 900, ""},
				{time.Date(2023, 3, 4, 0, 0, 0, 0, currentDate.Location()), 1800, 900, ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateOccurrences(tt.cycle, tt.given)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateOccurrences(%#v, %v)\n  got = %v,\n  want  %v", tt.cycle, tt.given, got, tt.want)
			}
		})
	}
}

func TestApplyExceptions(t *testing.T) {
	currentDate := time.Now()
	jan1 := time.Date(2023, 1, 1, 0, 0, 0, 0, currentDate.Location())
	dec31 := time.Date(2023, 12, 31, 0, 0, 0, 0, currentDate.Location())
	mar1 := time.Date(2023, 3, 1, 0, 0, 0, 0, currentDate.Location())
	mar2 := time.Date(2023, 3, 2, 0, 0, 0, 0, currentDate.Location())
	mar5 := time.Date(2023, 3, 5, 0, 0, 0, 0, currentDate.Location())
	mar31 := time.Date(2023, 3, 31, 0, 0, 0, 0, currentDate.Location())
	tests := []struct {
		name        string
		occurrences []domain.Occurrence
		exceptions  []domain.Exception
		start       time.Time
		end         time.Time
		want        []domain.Occurrence
	}{
		{"exception adds location",
			[]domain.Occurrence{
				{mar1, 1800, 900, ""},
				{mar2, 1800, 900, ""},
				{mar5, 1800, 900, ""},
			},
			[]domain.Exception{
				{mar2, 1800, 900, "123"},
			},
			mar1, mar31,
			[]domain.Occurrence{
				{mar1, 1800, 900, ""},
				{mar2, 1800, 900, "123"},
				{mar5, 1800, 900, ""},
			},
		},
		{"two events are cancelled",
			[]domain.Occurrence{
				{mar1, 1800, 900, ""},
				{mar2, 1800, 900, ""},
				{mar5, 1800, 900, ""},
			},
			[]domain.Exception{
				{mar2, 0, 0, ""},
				{mar5, 0, 0, ""},
			},
			mar1, mar31,
			[]domain.Occurrence{{mar1, 1800, 900, ""}},
		},
		{"exception adds one more day",
			[]domain.Occurrence{
				{mar1, 1800, 900, ""},
				{mar5, 1800, 900, ""},
			},
			[]domain.Exception{
				{mar2, 1800, 900, ""},
			},
			mar1, mar31,
			[]domain.Occurrence{
				{mar1, 1800, 900, ""},
				{mar2, 1800, 900, ""},
				{mar5, 1800, 900, ""},
			},
		},
		{"exception is before or after timespan",
			[]domain.Occurrence{
				{mar1, 1800, 900, ""},
				{mar2, 1800, 900, ""},
				{mar5, 1800, 900, ""},
			},
			[]domain.Exception{
				{jan1, 1800, 900, "123"},
				{dec31, 1800, 900, "123"},
			},
			mar1, mar31,
			[]domain.Occurrence{
				{mar1, 1800, 900, ""},
				{mar2, 1800, 900, ""},
				{mar5, 1800, 900, ""},
			},
		},
		{"exception reschedules one day into two events on the same day",
			[]domain.Occurrence{
				{mar1, 1800, 900, ""},
				{mar2, 1800, 900, ""},
				{mar5, 1800, 900, ""},
			},
			[]domain.Exception{
				{mar2, 2400, 450, ""},
				{mar2, 3600, 450, ""},
			},
			mar1, mar31,
			[]domain.Occurrence{
				{mar1, 1800, 900, ""},
				{mar2, 2400, 450, ""},
				{mar2, 3600, 450, ""},
				{mar5, 1800, 900, ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyExceptions(tt.occurrences, tt.exceptions, tt.start, tt.end)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplyExceptions(%v, %v, %v, %v)\n  got = %v,\n  want  %v", tt.occurrences, tt.exceptions, tt.start, tt.end, got, tt.want)
			}
		})
	}
}

func TestTrimOccurrences(t *testing.T) {
	currentDate := time.Now()
	mar1 := time.Date(2023, 3, 1, 0, 0, 0, 0, currentDate.Location())
	mar2 := time.Date(2023, 3, 2, 0, 0, 0, 0, currentDate.Location())
	mar5 := time.Date(2023, 3, 5, 0, 0, 0, 0, currentDate.Location())
	mar31 := time.Date(2023, 3, 31, 0, 0, 0, 0, currentDate.Location())
	tests := []struct {
		name        string
		occurrences []domain.Occurrence
		start       time.Time
		end         time.Time
		want        []domain.Occurrence
	}{
		{"trim one day",
			[]domain.Occurrence{
				{mar1, 1800, 900, ""},
				{mar2, 1800, 900, ""},
				{mar5, 1800, 900, ""},
			},
			mar2, mar31,
			[]domain.Occurrence{
				{mar2, 1800, 900, ""},
				{mar5, 1800, 900, ""},
			},
		},
		{"trim one day, short period",
			[]domain.Occurrence{
				{mar1, 1800, 900, ""},
				{mar2, 1800, 900, ""},
				{mar5, 1800, 900, ""},
			},
			mar2, mar2,
			[]domain.Occurrence{
				{mar2, 1800, 900, ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trimOccurrences(tt.occurrences, tt.start, tt.end)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trimOccurrences(%v, %v, %v)\n  got = %v,\n  want  %v", tt.occurrences, tt.start, tt.end, got, tt.want)
			}
		})
	}
}
