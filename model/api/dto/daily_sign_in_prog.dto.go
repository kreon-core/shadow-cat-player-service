package dto

type DailySignInProgress struct {
	WeekID      string       `json:"week_id"`
	ClaimedDays map[int]bool `json:"claimed_days"`
}
