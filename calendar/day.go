package calendar

import (
	"errors"
	"time"
)

var Weekdays = []time.Weekday{ time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday }
var WeekdayNamesDe = []string{ "Mo", "Di", "Mi", "Do", "Fr", "Sa", "So" }

var InvalidWeekday = errors.New("invalid weekday")

func WeekdayIndex(weekday time.Weekday) int {
	switch weekday {
	case time.Monday:
		return 0
	case time.Tuesday:
		return 1
	case time.Wednesday:
		return 2
	case time.Thursday:
		return 3
	case time.Friday:
		return 4
	case time.Saturday:
		return 5
	case time.Sunday:
		return 6
	}

	panic(InvalidWeekday)
}

func FormatWeekday(weekday time.Weekday) string {
	index := WeekdayIndex(weekday)
	return WeekdayNamesDe[index]
}

func FormatWeekdays(weekdays ...time.Weekday) []string {
	names := make([]string, 0)

	for _, weekday := range weekdays {
		names = append(names, FormatWeekday(weekday))
	}

	return names
}

func Day(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func Today() time.Time {
	return Day(time.Now())
}