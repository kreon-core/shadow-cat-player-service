package helper

import "time"

func WeekStartAt(now time.Time) time.Time {
	year, week := now.ISOWeek()

	loc := now.Location()
	t := time.Date(year, 1, 1, 0, 0, 0, 0, loc)

	isoStart := t.AddDate(0, 0, (week-1)*7)
	for isoStart.Weekday() != time.Monday {
		isoStart = isoStart.AddDate(0, 0, -1)
	}
	return isoStart
}

func DayNoFromStartWeek(now time.Time) int {
	return int(now.Weekday()+6) % 7
}
