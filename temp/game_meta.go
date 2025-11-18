package temp

import "time"

const (
	BasicCoins  = 20000
	BasicGems   = 100
	BasicEnergy = 30

	EnergyRegenInterval = 15 * time.Minute
)

var UnlockDailySignInCosts = []int{100, 500, 1000, 1500, 2000, 2500, 3000}
