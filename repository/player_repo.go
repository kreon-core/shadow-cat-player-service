package repository

import "sc-player-service/repository/pgsqlc"

type Player struct {
	q *pgsqlc.Queries
}

func NewPlayer(q *pgsqlc.Queries) *Player {
	return &Player{q: q}
}
