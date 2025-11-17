package response

type TowerProgress struct {
	Tower []Tower `json:"played_towers"`
}

type Tower struct {
	TowerID      int `json:"tower_id"`
	Ticket       int `json:"ticket"`
	HighestFloor int `json:"highest_floor"`
}
