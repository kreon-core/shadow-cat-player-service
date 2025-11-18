package request

type UnlockDailySignIn struct {
	WeekID string `json:"week_id" binding:"required"`
	DayNo  int    `json:"day_no"  binding:"required,min=0,max=6"`
}
