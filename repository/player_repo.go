package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"sc-player-service/repository/playersqlc"
)

type Player struct {
	q *playersqlc.Queries
}

func NewPlayer(q *playersqlc.Queries) *Player {
	return &Player{q: q}
}

func (r *Player) GetPlayerByID(ctx context.Context, id pgtype.UUID) (*playersqlc.GetPlayerByIDRow, error) {
	res, err := r.q.GetPlayerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get_player_by_id -> %w", err)
	}
	return &res, nil
}

func (r *Player) CreatePlayer(ctx context.Context, player *playersqlc.CreatePlayerParams) error {
	_, err := r.q.CreatePlayer(ctx, *player)
	if err != nil {
		return fmt.Errorf("create_player -> %w", err)
	}
	return nil
}
