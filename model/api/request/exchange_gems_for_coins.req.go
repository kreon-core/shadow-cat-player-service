package request

type ExchangeGemsForCoins struct {
	GemsCost    int32 `json:"gems_cost"    binding:"required,min=1"`
	CoinsGained int32 `json:"coins_gained" binding:"required,min=1"`
}
