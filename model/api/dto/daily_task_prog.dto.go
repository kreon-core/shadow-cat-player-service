package response

type DailyTaskProgress struct {
	Tasks []DailyTask `json:"tasks"`
}

type DailyTask struct {
	TaskID       string `json:"task_id"`
	Progress     int    `json:"progress"`
	Claimed      bool   `json:"claimed"`
	PointsEarned int    `json:"points_earned"`
}
