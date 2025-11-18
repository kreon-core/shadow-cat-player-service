package dto

type DailyTaskProgress struct {
	TotalPoints   int32          `json:"total_points"`
	Tasks         []DailyTask    `json:"tasks"`
	PlayerChanges *PlayerChanges `json:"player_changes,omitempty"`
}

type DailyTask struct {
	TaskID       int32 `json:"task_id"`
	Progress     int32 `json:"progress"`
	Claimed      bool  `json:"claimed"`
	PointsEarned int32 `json:"points_earned"`
}
