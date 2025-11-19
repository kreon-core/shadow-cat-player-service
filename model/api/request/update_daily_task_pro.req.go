package request

type UpdateDailyTaskProgress struct {
	TaskChanges []TaskChange `json:"task_changes" validate:"required,dive"`
}

type TaskChange struct {
	TaskID   int32 `json:"task_id"  validate:"required"`
	Progress int32 `json:"progress" validate:"required"`
}
