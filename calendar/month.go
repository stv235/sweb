package calendar

import "time"

var MonthNamesDe = []string{ "Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember" }

func FormatMonthName(month time.Month) string {
	return MonthNamesDe[int(month) - 1]
}

func Month(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}

func ThisMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}

func NextMonth(t time.Time) time.Time {
	year := t.Year()
	month := t.Month()

	month++

	if month > 12 {
		month = 1
		year++
	}

	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func LastMonth(t time.Time) time.Time {
	year := t.Year()
	month := t.Month()

	month--

	if month < 1 {
		month = 12
		year++
	}

	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}