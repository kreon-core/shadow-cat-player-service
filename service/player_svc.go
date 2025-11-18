package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kreon-core/shadow-cat-common/dbc"

	"sc-player-service/model/api/dto"
	"sc-player-service/model/api/request"
	"sc-player-service/repository"
	"sc-player-service/repository/playersqlc"
	"sc-player-service/temp"
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
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	// TODO: Begin transaction?

	createPlayerIfNotExists := false
	player, err := s.PlayerRepo.PlayerQueries.GetPlayerByID(ctx, id)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get_player_by_id -> %w", err)
		}
		createPlayerIfNotExists = true
	}

	if createPlayerIfNotExists {
		newPlayer, err := s.newPlayer(id)
		if err != nil {
			return nil, fmt.Errorf("create_new_player_struct -> %w", err)
		}

		_, err = s.PlayerRepo.PlayerQueries.CreateNewPlayer(ctx, *newPlayer)
		if err != nil {
			return nil, fmt.Errorf("create_player -> %w", err)
		}

		player, err = s.PlayerRepo.PlayerQueries.GetPlayerByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("get_player_by_id_after_create -> %w", err)
		}
	}

	// TODO: End transaction

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

	return &dto.Player{
		PlayerID: player.ID.String(),
		Level:    int(player.Level),
		EXP:      int(player.Exp),
		Coins:    int(player.Coins),
		Gems:     int(player.Gems),

		BestMap: bestMap,

		CurrentSkin:   int(player.CurrentSkin),
		EquippedProps: equippedProps,
	}, nil
}

func (s *Player) UpdatePlayer(
	ctx context.Context,
	playerID string,
	updateData *request.UpdatePlayer,
) (*dto.Player, error) {
	_, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	return nil, nil
}

func (s *Player) newPlayer(id pgtype.UUID) (*playersqlc.CreateNewPlayerParams, error) {
	bestMap := dto.BestMap{
		MapID:      0,
		TimeRecord: "00:00:00",
	}
	bestMapBytes, err := json.Marshal(bestMap)
	if err != nil {
		return nil, fmt.Errorf("marshal_best_map -> %w", err)
	}

	return &playersqlc.CreateNewPlayerParams{
		ID:            id,
		Coins:         temp.BasicCoins,
		Gems:          temp.BasicGems,
		CurrentEnergy: temp.BasicEnergy,
		MaxEnergy:     temp.BasicEnergy,
		BestMap:       bestMapBytes,
	}, nil
}
