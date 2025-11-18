package request

type ExchangeGemsForCoins struct {
	GemsCost    int32 `json:"gems_cost"    validate:"required,min=1"`
	CoinsGained int32 `json:"coins_gained" validate:"required,min=1"`
}
