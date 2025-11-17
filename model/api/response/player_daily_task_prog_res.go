package response

type PlayerDailyTaskProgress struct {
	Tasks []PlayerDailyTask `json:"tasks"`
}

type PlayerDailyTask struct {
	TaskID       string `json:"task_id"`
	Progress     int    `json:"progress"`
	Claimed      bool   `json:"claimed"`
	PointsEarned int    `json:"points_earned"`
}
