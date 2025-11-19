package request

type ClaimDailyTask struct {
	TaskID       int32 `json:"task_id"       validate:"required"`
	PointsEarned int32 `json:"points_earned" validate:"required,min=0"`
}
