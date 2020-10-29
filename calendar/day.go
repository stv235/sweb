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

func daysToNext(w1, w2 time.Weekday) int {
	d1 := int(w1)
	d2 := int(w2)

	if d1 == d2 {
		return 7
	}

	if d1 < d2 {
		return d2 - d1
	}

	return d2 + 7 - d1
}

func NextWeekday(t time.Time, weekday time.Weekday) time.Time {
	return t.AddDate(0, 0, daysToNext(t.Weekday(), weekday))
}

func PreviousMonday(t time.Time) time.Time {
	for t.Weekday() != time.Monday {
		t = t.AddDate(0, 0, -1)
	}

	return t
}