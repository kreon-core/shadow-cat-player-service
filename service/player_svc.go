package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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

	return s.getPlayerData(&player)
}

func (s *Player) UpdatePlayer(
	ctx context.Context, playerID string,
	updateData *request.UpdatePlayer,
) (*dto.Player, error) {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	// TODO: Begin transaction?

	player, err := s.PlayerRepo.PlayerQueries.GetPlayerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get_player_by_id -> %w", err)
	}

	equippedPropsBytes, err := json.Marshal(updateData.EquippedProps)
	if err != nil {
		return nil, fmt.Errorf("marshal_best_map -> %w", err)
	}
	_, err = s.PlayerRepo.PlayerQueries.UpdatePlayer(ctx, playersqlc.UpdatePlayerParams{
		ID:            player.ID,
		Level:         player.Level,
		Exp:           player.Exp,
		Coins:         player.Coins,
		Gems:          player.Gems,
		BestMap:       player.BestMap,
		CurrentSkin:   updateData.CurrentSkin,
		EquippedProps: equippedPropsBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("update_player -> %w", err)
	}

	// TODO: End transaction

	return s.getPlayerData(&player)
}

func (s *Player) GetEnergy(ctx context.Context, playerID string) (*dto.PlayerEnergy, error) {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	// TODO: Begin transaction?

	playerEnergy, err := s.PlayerRepo.PlayerQueries.GetPlayerEnergyByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get_player_energy_by_id -> %w", err)
	}

	maxEnergy := playerEnergy.MaxEnergy
	curEnergy, nextEnergyAt, err := s.calcNUpdateEnergy(ctx, id, &playerEnergy)
	if err != nil {
		return nil, fmt.Errorf("calc_n_update_energy -> %w", err)
	}

	// TODO: End transaction

	nextEnergyAtUnix := int64(-1)
	if nextEnergyAt.Valid {
		nextEnergyAtUnix = nextEnergyAt.Time.Unix()
	}
	return &dto.PlayerEnergy{
		CurrentEnergy: int(curEnergy),
		MaxEnergy:     int(maxEnergy),
		NextEnergyAt:  nextEnergyAtUnix,
	}, nil
}

func (s *Player) GetInventory(ctx context.Context, playerID string) (*dto.PlayerInventory, error) {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	playerInventory, err := s.PlayerRepo.PlayerQueries.GetInventoryByPlayerID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get_inventory_by_player_id -> %w", err)
	}

	var skins []int
	switch v := playerInventory.OwnedSkins.(type) {
	case []int64:
		for _, s := range v {
			skins = append(skins, int(s))
		}
	case nil:
		skins = []int{}
	default:
		return nil, fmt.Errorf("unexpected type for OwnedSkins: %T", v)
	}

	var props []dto.Prop
	switch v := playerInventory.OwnedProps.(type) {
	case []byte:
		if err := json.Unmarshal(v, &props); err != nil {
			return nil, fmt.Errorf("failed to unmarshal props: %w", err)
		}
	case nil:
		props = []dto.Prop{}
	default:
		return nil, fmt.Errorf("unexpected type for OwnedProps: %T", v)
	}

	return &dto.PlayerInventory{
		Skins: skins,
		Props: props,
	}, nil
}

func (s *Player) GetTowerProgress(ctx context.Context, playerID string) (*dto.TowerProgress, error) {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	towerProgress, err := s.PlayerRepo.PlayerQueries.GetTowerProgressByPlayerID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get_tower_progress_by_player_id -> %w", err)
	}

	towerProgressData := &dto.TowerProgress{}
	for _, tp := range towerProgress {
		towerProgressData.Tower = append(towerProgressData.Tower, dto.Tower{
			TowerID:      int(tp.TowerID),
			Ticket:       int(tp.Ticket),
			HighestFloor: int(tp.HighestFloor),
		})
	}

	return towerProgressData, nil
}

func (s *Player) GetChapterProgress(ctx context.Context, playerID string) (*dto.ChapterProgress, error) {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	chapterProgress, err := s.PlayerRepo.PlayerQueries.GetChapterProgressByPlayerID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get_chapter_progress_by_player_id -> %w", err)
	}

	chapterProgressData := &dto.ChapterProgress{}
	for _, cp := range chapterProgress {
		checkedCheckpoints := make(map[int]bool)
		err := json.Unmarshal(cp.CheckedCheckpoints, &checkedCheckpoints)
		if err != nil {
			return nil, fmt.Errorf("unmarshal_checked_checkpoints -> %w", err)
		}
		chapterProgressData.Chapters = append(chapterProgressData.Chapters, dto.Chapter{
			ChapterID:          int(cp.ChapterID),
			CheckedCheckpoints: checkedCheckpoints,
		})
	}

	return chapterProgressData, nil
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

func (s *Player) getPlayerData(player *playersqlc.GetPlayerByIDRow) (*dto.Player, error) {
	var bestMap dto.BestMap
	err := json.Unmarshal(player.BestMap, &bestMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshal_best_map -> %w", err)
	}

	var equippedProps []int
	err = json.Unmarshal(player.EquippedProps, &equippedProps)
	if err != nil {
		return nil, fmt.Errorf("unmarshal_equipped_props -> %w", err)
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

func (s *Player) calcNUpdateEnergy(
	ctx context.Context,
	id pgtype.UUID,
	playerEnergy *playersqlc.GetPlayerEnergyByIDRow,
) (int32, pgtype.Timestamptz, error) {
	curEnergy := playerEnergy.CurrentEnergy
	maxEnergy := playerEnergy.MaxEnergy
	nextEnergyAt := playerEnergy.NextEnergyAt
	if curEnergy >= maxEnergy || !nextEnergyAt.Valid {
		return curEnergy, nextEnergyAt, nil
	}

	now := time.Now()
	if now.Before(nextEnergyAt.Time) {
		return curEnergy, nextEnergyAt, nil
	}

	elapsed := now.Sub(nextEnergyAt.Time)
	energyToRegen := int32(elapsed / temp.EnergyRegenInterval) //nolint:gosec // integer division
	if energyToRegen > 0 {
		newEnergy := min(curEnergy+energyToRegen, maxEnergy)

		var newNextEnergyAt pgtype.Timestamptz
		if newEnergy < maxEnergy {
			remainder := elapsed % temp.EnergyRegenInterval
			newNextEnergyAt = pgtype.Timestamptz{
				Time:  now.Add(temp.EnergyRegenInterval - remainder),
				Valid: true,
			}
		} else {
			newNextEnergyAt = pgtype.Timestamptz{Valid: false}
		}

		_, err := s.PlayerRepo.PlayerQueries.UpdatePlayerEnergy(ctx, playersqlc.UpdatePlayerEnergyParams{
			ID:            id,
			CurrentEnergy: newEnergy,
			NextEnergyAt:  newNextEnergyAt,
		})
		if err != nil {
			return 0, pgtype.Timestamptz{}, fmt.Errorf("update_player_energy -> %w", err)
		}

		curEnergy = newEnergy
		nextEnergyAt = newNextEnergyAt
	}
	return curEnergy, nextEnergyAt, nil
}
