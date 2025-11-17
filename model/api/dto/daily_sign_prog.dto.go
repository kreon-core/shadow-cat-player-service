package response

type DailySignProgress struct {
	WeekNo     int         `json:"week_no"`
	SignedDays []SignedDay `json:"signed_days"`
}

type SignedDay struct {
	Day     int  `json:"day"`
	Claimed bool `json:"claimed"`
}
