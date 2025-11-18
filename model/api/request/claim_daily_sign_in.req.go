package request

type ClaimDailySignIn struct {
	WeekID string `json:"week_id" validate:"required"`
	DayNo  int    `json:"day_no"  validate:"required,min=0,max=6"`
}
