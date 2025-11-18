package helper

import "time"

func WeekStartAt() time.Time {
	now := time.Now()
	year, week := now.ISOWeek()

	loc := now.Location()
	t := time.Date(year, 1, 1, 0, 0, 0, 0, loc)

	isoStart := t.AddDate(0, 0, (week-1)*7)
	for isoStart.Weekday() != time.Monday {
		isoStart = isoStart.AddDate(0, 0, -1)
	}
	return isoStart
}
