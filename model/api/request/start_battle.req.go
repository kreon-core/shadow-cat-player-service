package request

type StartBattle struct {
	GameMode int32  `json:"game_mode" validate:"required"`
	TowerID  *int32 `json:"tower_id"`
	Floor    *int32 `json:"floor"`
	MapID    *int32 `json:"map_id"`
}
