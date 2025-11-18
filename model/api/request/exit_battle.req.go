package request

type ExitBattle struct {
	BattleID string `json:"battle_id" validate:"required"`
}
