package repository

import "sc-player-service/repository/playersqlc"

type Player struct {
	q *playersqlc.Queries
}

func NewPlayer(q *playersqlc.Queries) *Player {
	return &Player{q: q}
}
