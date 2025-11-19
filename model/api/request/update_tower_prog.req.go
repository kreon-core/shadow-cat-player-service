package request

type UpdateTowerProgress struct {
	TowerID      int32 `json:"tower_id"      validate:"required"`
	HighestFloor int32 `json:"highest_floor" validate:"required,min=0"`
}
