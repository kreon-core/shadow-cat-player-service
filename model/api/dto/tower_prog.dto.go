package dto

type TowerProgress struct {
	PlayedTower []Tower `json:"played_towers"`
}

type Tower struct {
	TowerID      int32 `json:"tower_id"`
	Ticket       int32 `json:"ticket"`
	HighestFloor int32 `json:"highest_floor"`
}
