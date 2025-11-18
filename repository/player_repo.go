package repository

import (
	"sc-player-service/repository/playersqlc"
)

type Player struct {
	PlayerQueries *playersqlc.Queries
}

func NewPlayer(playerQueries *playersqlc.Queries) *Player {
	return &Player{PlayerQueries: playerQueries}
}
