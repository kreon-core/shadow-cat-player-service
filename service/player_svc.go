package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"sc-player-service/model/api/dto"
	"sc-player-service/repository"
	"sc-player-service/repository/playersqlc"
)

type Player struct {
	PlayerRepo *repository.Player
}

func NewPlayer(playerRepo *repository.Player) *Player {
	return &Player{
		PlayerRepo: playerRepo,
	}
}

func (s *Player) GetOrCreatePlayer(ctx context.Context, playerID string) (*dto.Player, error) {
	var id pgtype.UUID
	err := id.Scan(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	createPlayerIfNotExists := false
	player, err := s.PlayerRepo.GetPlayerByID(ctx, id)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("player_repo_get_by_id -> %w", err)
		}
		createPlayerIfNotExists = true
	}

	if createPlayerIfNotExists {
		bestMap := dto.BestMap{
			MapID:      0,
			TimeRecord: "00:00:00",
		}
		bestMapBytes, err := json.Marshal(bestMap)
		if err != nil {
			return nil, fmt.Errorf("marshal_best_map -> %w", err)
		}

		newPlayer := &playersqlc.CreatePlayerParams{
			ID:      id,
			Coins:   100,
			Gems:    10,
			BestMap: bestMapBytes,
		}

		err = s.PlayerRepo.CreatePlayer(ctx, newPlayer)
		if err != nil {
			return nil, fmt.Errorf("player_repo_create -> %w", err)
		}

		player, err = s.PlayerRepo.GetPlayerByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("player_repo_get_by_id_after_create -> %w", err)
		}
	}

	var bestMap dto.BestMap
	err = json.Unmarshal(player.BestMap, &bestMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshal_best_map -> %w", err)
	}

	var equippedProps []int
	if len(player.EquippedProps) > 0 {
		if err := json.Unmarshal(player.EquippedProps, &equippedProps); err != nil {
			return nil, fmt.Errorf("unmarshal_equipped_props -> %w", err)
		}
	}

	result := &dto.Player{
		PlayerID: player.ID.String(),
		Level:    int(player.Level),
		EXP:      int(player.Exp),
		Coins:    int(player.Coins),
		Gems:     int(player.Gems),

		BestMap: bestMap,

		CurrentSkin:   int(player.CurrentSkin),
		EquippedProps: equippedProps,
	}

	return result, nil
}
