package service

import "sc-player-service/repository"

type Player struct {
	PlayerRepo *repository.Player
}

func NewPlayer(playerRepo *repository.Player) *Player {
	return &Player{
		PlayerRepo: playerRepo,
	}
}
