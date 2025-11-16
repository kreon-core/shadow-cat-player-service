package service

import repository "sc-player-service/repository/postgres"

type Player struct {
	PlayerRepo *repository.Player
}

func NewPlayer(playerRepo *repository.Player) *Player {
	return &Player{
		PlayerRepo: playerRepo,
	}
}
