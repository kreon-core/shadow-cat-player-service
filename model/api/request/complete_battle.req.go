package request

import "sc-player-service/model/api/dto"

type CompleteBattle struct {
	GameMode     int32        `json:"game_mode"     validate:"required"`
	TowerID      *int32       `json:"tower_id"`
	Floor        *int32       `json:"floor"`
	MapID        *int32       `json:"map_id"`
	BattleResult BattleResult `json:"battle_result" validate:"required,dive"`
}

type BattleResult struct {
	HeroEXP          int32      `json:"hero_exp"           validate:"required,min=0"`
	TimeSurvived     int32      `json:"time_survived"      validate:"required"`
	MonsterKills     int32      `json:"monster_kills"      validate:"required,min=0"`
	TotalDamageDealt int32      `json:"total_damage_dealt" validate:"required,min=0"`
	CoinsCollected   int32      `json:"coins_collected"    validate:"required,min=0"`
	GemsCollected    int32      `json:"gems_collected"     validate:"required,min=0"`
	Props            []dto.Prop `json:"props"              validate:"dive"`
}
