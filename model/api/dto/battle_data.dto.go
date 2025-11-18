package dto

type BattleData struct {
	BattleID      string         `json:"battle_id"`
	PlayerChanges *PlayerChanges `json:"player_changes,omitempty"`
}
