package response

type PlayerTowerProgress struct {
	PlayedTowers []PlayerTower `json:"played_towers"`
}

type PlayerTower struct {
	TowerID   int `json:"tower_id"`
	Ticket    int `json:"ticket"`
	BestFloor int `json:"best_floor"`
}
