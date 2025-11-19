package request

type UpdateTowerProgress struct {
	TowerID      int32 `json:"tower_id"`
	HighestFloor int32 `json:"highest_floor"`
}
