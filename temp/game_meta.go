package temp

import (
	"time"

	"sc-player-service/model/api/dto"
)

const (
	BasicCoins  = 20000
	BasicGems   = 100
	BasicEnergy = 30

	EnergyRegenInterval = 15 * time.Minute
)

type Reward struct {
	Coins int32      `json:"coins,omitempty"`
	Gems  int32      `json:"gems,omitempty"`
	Props []dto.Prop `json:"props,omitempty"`
}

var (
	UnlockDailySignInCosts = []int32{100, 500, 1000, 1500, 2000, 2500, 3000}
	DailySignInRewards     = []Reward{
		{Coins: 1000, Gems: 10},
		{Coins: 2000, Gems: 15},
		{Coins: 3000, Gems: 20},
		{Coins: 4000, Gems: 25},
		{Coins: 5000, Gems: 30},
		{Coins: 6000, Gems: 35},
		{Coins: 10000, Gems: 50},
	}

	ChapterCheckpointRewards = map[int32]Reward{
		1: {Coins: 500, Gems: 5},
		2: {Coins: 1000, Gems: 10},
		3: {Coins: 1500, Gems: 15},
	}
)
