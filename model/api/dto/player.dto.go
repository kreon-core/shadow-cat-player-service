package dto

type Player struct {
	PlayerID string `json:"player_id"`
	Level    int32  `json:"level"`
	EXP      int32  `json:"exp"`
	Coins    int32  `json:"coins"`
	Gems     int32  `json:"gems"`

	BestMap BestMap `json:"best_map"`

	CurrentSkin   int32   `json:"current_skin"`
	EquippedProps []int32 `json:"equipped_props"`
}

type BestMap struct {
	MapID      int32  `json:"map_id"`
	TimeRecord string `json:"time_record"`
}
