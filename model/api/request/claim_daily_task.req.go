package request

type ClaimDailyTask struct {
	TaskID       int32 `json:"task_id"       binding:"required"`
	PointsEarned int32 `json:"points_earned" binding:"required,min=1"`
}
