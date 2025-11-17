package response

type Player struct {
	PlayerID string `json:"player_id"`
	Level    int    `json:"level"`
	EXP      int    `json:"exp"`
	Coins    int    `json:"coins"`
	Gems     int    `json:"gems"`

	BestMap BestMap `json:"best_map"`

	CurrentSkin   int   `json:"current_skin"`
	EquippedProps []int `json:"equipped_props"`
}

type BestMap struct {
	MapID      int    `json:"map_id"`
	TimeRecord string `json:"time_record"`
}
