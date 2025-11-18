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

	"sc-player-service/helper"
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

	player, err := s.updatePlayer(ctx, id, &dto.PlayerChanges{
		CurrentSkin:   &updateData.CurrentSkin,
		EquippedProps: &updateData.EquippedProps,
	})
	if err != nil {
		return nil, fmt.Errorf("update_player -> %w", err)
	}

	return player, nil
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

func (s *Player) ClaimChapterRewards(
	ctx context.Context,
	playerID string,
	req *request.ClaimChapterRewards,
) (*dto.PlayerChanges, error) {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	chapter, err := s.PlayerRepo.PlayerQueries.GetChapterProgressByPlayerIDAndChapterID(
		ctx,
		playersqlc.GetChapterProgressByPlayerIDAndChapterIDParams{
			PlayerID:  id,
			ChapterID: req.ChapterID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get_chapter_progress_by_player_id_and_chapter_id -> %w", err)
	}

	checkedCheckpoints := make(map[int32]bool)
	err = json.Unmarshal(chapter.CheckedCheckpoints, &checkedCheckpoints)
	if err != nil {
		return nil, fmt.Errorf("unmarshal_checked_checkpoints -> %w", err)
	}

	if v, ok := checkedCheckpoints[req.Checkpoint]; ok && v {
		return nil, errors.New("checkpoint_already_claimed")
	}
	checkedCheckpoints[req.Checkpoint] = true
	checkedCheckpointsBytes, err := json.Marshal(checkedCheckpoints)
	if err != nil {
		return nil, fmt.Errorf("marshal_checked_checkpoints -> %w", err)
	}

	_, err = s.PlayerRepo.PlayerQueries.UpsertChapterProgressOnPlayer(
		ctx,
		playersqlc.UpsertChapterProgressOnPlayerParams{
			PlayerID:           id,
			ChapterID:          req.ChapterID,
			CheckedCheckpoints: checkedCheckpointsBytes,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("update_chapter_checked_checkpoints -> %w", err)
	}

	reward := temp.ChapterCheckpointRewards[req.Checkpoint]
	player, err := s.updatePlayer(ctx, id, &dto.PlayerChanges{
		Coins: &reward.Coins,
		Gems:  &reward.Gems,
	})
	if err != nil {
		return nil, fmt.Errorf("update_player -> %w", err)
	}

	return &dto.PlayerChanges{
		Coins: &player.Coins,
		Gems:  &player.Gems,
	}, nil
}

func (s *Player) GetDailySignInProgress(ctx context.Context, playerID string) (*dto.DailySignInProgress, error) {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	markDailySignInDaysParams, err := s.getOrCreateDailySignInProgress(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get_or_create_daily_sign_in_progress -> %w", err)
	}

	claimedDays := make(map[int]bool)
	err = json.Unmarshal(markDailySignInDaysParams.ClaimedDays, &claimedDays)
	if err != nil {
		return nil, fmt.Errorf("unmarshal_claimed_days -> %w", err)
	}

	return &dto.DailySignInProgress{
		WeekID:      markDailySignInDaysParams.ID.String(),
		ClaimedDays: claimedDays,
	}, nil
}

func (s *Player) MarkDailySignIn(ctx context.Context, playerID string) error {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return fmt.Errorf("parse_uuid_string -> %w", err)
	}

	// TODO: Begin transaction?

	markDailySignInDaysParams, err := s.getOrCreateDailySignInProgress(ctx, id)
	if err != nil {
		return fmt.Errorf("get_or_create_daily_sign_in_progress -> %w", err)
	}

	claimedDays := make(map[int]bool)
	err = json.Unmarshal(markDailySignInDaysParams.ClaimedDays, &claimedDays)
	if err != nil {
		return fmt.Errorf("unmarshal_claimed_days -> %w", err)
	}

	dayNo := helper.DayNoFromStartWeek(time.Now())
	if v, ok := claimedDays[dayNo]; ok && v {
		return nil
	}
	claimedDays[dayNo] = false

	claimedDaysBytes, err := json.Marshal(claimedDays)
	if err != nil {
		return fmt.Errorf("marshal_claimed_days -> %w", err)
	}

	_, err = s.PlayerRepo.PlayerQueries.MarkDailySignInDays(ctx, playersqlc.MarkDailySignInDaysParams{
		ID:          markDailySignInDaysParams.ID,
		PlayerID:    markDailySignInDaysParams.PlayerID,
		ClaimedDays: claimedDaysBytes,
	})
	if err != nil {
		return fmt.Errorf("mark_daily_sign_in_days -> %w", err)
	}

	// TODO: End transaction

	return nil
}

func (s *Player) UnlockDailySignIn(
	ctx context.Context,
	playerID string,
	req *request.UnlockDailySignIn,
) (*dto.PlayerChanges, error) {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	weekID, err := dbc.ParseUUID(req.WeekID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string_week_id -> %w", err)
	}

	// TODO: Begin transaction?

	getDailySignInByIDRow, err := s.PlayerRepo.PlayerQueries.GetDailySignInByID(
		ctx,
		playersqlc.GetDailySignInByIDParams{
			ID:       weekID,
			PlayerID: id,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get_daily_sign_in_by_player_id -> %w", err)
	}

	claimedDays := make(map[int]bool)
	err = json.Unmarshal(getDailySignInByIDRow.ClaimedDays, &claimedDays)
	if err != nil {
		return nil, fmt.Errorf("unmarshal_claimed_days -> %w", err)
	}

	claimedDays[req.DayNo] = false
	claimedDaysBytes, err := json.Marshal(claimedDays)
	if err != nil {
		return nil, fmt.Errorf("marshal_claimed_days -> %w", err)
	}

	_, err = s.PlayerRepo.PlayerQueries.MarkDailySignInDays(ctx, playersqlc.MarkDailySignInDaysParams{
		ID:          getDailySignInByIDRow.ID,
		PlayerID:    id,
		ClaimedDays: claimedDaysBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("mark_daily_sign_in_days -> %w", err)
	}

	coinsCost := -temp.UnlockDailySignInCoinsCost[req.DayNo]
	player, err := s.updatePlayer(ctx, id, &dto.PlayerChanges{
		Coins: &coinsCost,
	})
	if err != nil {
		return nil, fmt.Errorf("update_player -> %w", err)
	}

	// TODO: End transaction

	return &dto.PlayerChanges{
		Coins: &player.Coins,
		Gems:  &player.Gems,
	}, nil
}

func (s *Player) ClaimDailySignIn(ctx context.Context,
	playerID string,
	req *request.ClaimDailySignIn,
) (*dto.PlayerChanges, error) {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	weekID, err := dbc.ParseUUID(req.WeekID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string_week_id -> %w", err)
	}

	// TODO: Begin transaction?

	getDailySignInByIDRow, err := s.PlayerRepo.PlayerQueries.GetDailySignInByID(
		ctx,
		playersqlc.GetDailySignInByIDParams{
			ID:       weekID,
			PlayerID: id,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get_daily_sign_in_by_player_id -> %w", err)
	}

	claimedDays := make(map[int]bool)
	err = json.Unmarshal(getDailySignInByIDRow.ClaimedDays, &claimedDays)
	if err != nil {
		return nil, fmt.Errorf("unmarshal_claimed_days -> %w", err)
	}

	claimedDays[req.DayNo] = true
	claimedDaysBytes, err := json.Marshal(claimedDays)
	if err != nil {
		return nil, fmt.Errorf("marshal_claimed_days -> %w", err)
	}

	_, err = s.PlayerRepo.PlayerQueries.MarkDailySignInDays(ctx, playersqlc.MarkDailySignInDaysParams{
		ID:          getDailySignInByIDRow.ID,
		PlayerID:    id,
		ClaimedDays: claimedDaysBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("mark_daily_sign_in_days -> %w", err)
	}

	reward := temp.DailySignInRewards[req.DayNo]
	player, err := s.updatePlayer(ctx, id, &dto.PlayerChanges{
		Coins: &reward.Coins,
		Gems:  &reward.Gems,
	})
	if err != nil {
		return nil, fmt.Errorf("update_player -> %w", err)
	}

	// TODO: End transaction

	return &dto.PlayerChanges{
		Coins: &player.Coins,
		Gems:  &player.Gems,
	}, nil
}

func (s *Player) GetDailyTaskProgress(ctx context.Context, playerID string) (*dto.DailyTaskProgress, error) {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	dayStartAt := helper.DayStartAt(time.Now())

	dailyTaskProgress, err := s.PlayerRepo.PlayerQueries.GetDailyTasksByPlayerID(
		ctx,
		playersqlc.GetDailyTasksByPlayerIDParams{
			PlayerID: id,
			DayStartAt: pgtype.Timestamptz{
				Time:  dayStartAt,
				Valid: true,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get_daily_task_progress_by_player_id -> %w", err)
	}

	dailyTaskProgressData := &dto.DailyTaskProgress{
		TotalPoints: 0,
	}
	for _, dtp := range dailyTaskProgress {
		dailyTaskProgressData.Tasks = append(dailyTaskProgressData.Tasks, dto.DailyTask{
			TaskID:       dtp.TaskID,
			Progress:     dtp.Progress,
			Claimed:      dtp.Claimed,
			PointsEarned: dtp.PointsEarned,
		})
		dailyTaskProgressData.TotalPoints += dtp.PointsEarned
	}

	return dailyTaskProgressData, nil
}

func (s *Player) ClaimDailyTask(
	ctx context.Context,
	playerID string,
	req *request.ClaimDailyTask,
) (*dto.DailyTaskProgress, error) {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	dayStartAt := helper.DayStartAt(time.Now())

	_, err = s.PlayerRepo.PlayerQueries.ClaimDailyTask(
		ctx,
		playersqlc.ClaimDailyTaskParams{
			PlayerID: id,
			TaskID:   req.TaskID,
			DayStartAt: pgtype.Timestamptz{
				Time:  dayStartAt,
				Valid: true,
			},
			PointsEarned: req.PointsEarned,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("claim_daily_task -> %w", err)
	}

	reward := temp.DailyTaskRewards[req.TaskID]
	player, err := s.updatePlayer(ctx, id, &dto.PlayerChanges{
		Coins: &reward.Coins,
		Gems:  &reward.Gems,
	})
	if err != nil {
		return nil, fmt.Errorf("update_player -> %w", err)
	}

	dailyTaskProgress, err := s.GetDailyTaskProgress(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("get_daily_task_progress -> %w", err)
	}

	dailyTaskProgress.PlayerChanges = &dto.PlayerChanges{
		Coins: &player.Coins,
		Gems:  &player.Gems,
	}

	return dailyTaskProgress, nil
}

func (s *Player) ExchangeGemsForCoins(
	ctx context.Context,
	playerID string,
	req *request.ExchangeGemsForCoins,
) (*dto.PlayerChanges, error) {
	id, err := dbc.ParseUUID(playerID)
	if err != nil {
		return nil, fmt.Errorf("parse_uuid_string -> %w", err)
	}

	gemsChange := -req.GemsCost
	player, err := s.updatePlayer(ctx, id, &dto.PlayerChanges{
		Gems:  &gemsChange,
		Coins: &req.CoinsGained,
	})
	if err != nil {
		return nil, fmt.Errorf("update_player -> %w", err)
	}

	return &dto.PlayerChanges{
		Coins: &player.Coins,
		Gems:  &player.Gems,
	}, nil
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

	var equippedProps []int32
	err = json.Unmarshal(player.EquippedProps, &equippedProps)
	if err != nil {
		return nil, fmt.Errorf("unmarshal_equipped_props -> %w", err)
	}

	return &dto.Player{
		PlayerID: player.ID.String(),
		Level:    player.Level,
		EXP:      player.Exp,
		Coins:    player.Coins,
		Gems:     player.Gems,

		BestMap: bestMap,

		CurrentSkin:   player.CurrentSkin,
		EquippedProps: equippedProps,
	}, nil
}

func (s *Player) updatePlayer(ctx context.Context, id pgtype.UUID, changes *dto.PlayerChanges) (*dto.Player, error) {
	// TODO: Begin transaction?

	player, err := s.PlayerRepo.PlayerQueries.GetPlayerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get_player_by_id -> %w", err)
	}

	if changes.Level != nil {
		player.Level += *changes.Level
		if player.Level < 0 {
			return nil, errors.New("insufficient_level")
		}
	}
	if changes.Exp != nil {
		player.Exp += *changes.Exp
		if player.Exp < 0 {
			return nil, errors.New("insufficient_exp")
		}

		// TODO: make this configurable
		if player.Exp > 100 {
			player.Level += player.Exp / 100
			player.Exp = player.Exp % 100
		}
	}
	if changes.Coins != nil {
		player.Coins += *changes.Coins
		if player.Coins < 0 {
			return nil, errors.New("insufficient_coins")
		}
	}
	if changes.Gems != nil {
		player.Gems += *changes.Gems
		if player.Gems < 0 {
			return nil, errors.New("insufficient_gems")
		}
	}
	if changes.BestMap != nil {
		bestMapBytes, err := json.Marshal(changes.BestMap)
		if err != nil {
			return nil, fmt.Errorf("marshal_best_map -> %w", err)
		}
		player.BestMap = bestMapBytes
	}
	if changes.CurrentSkin != nil {
		player.CurrentSkin = *changes.CurrentSkin
	}
	if changes.EquippedProps != nil {
		equippedPropsBytes, err := json.Marshal(changes.EquippedProps)
		if err != nil {
			return nil, fmt.Errorf("marshal_equipped_props -> %w", err)
		}
		player.EquippedProps = equippedPropsBytes
	}

	params := playersqlc.UpdatePlayerParams(player)
	_, err = s.PlayerRepo.PlayerQueries.UpdatePlayer(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("update_player -> %w", err)
	}

	player, err = s.PlayerRepo.PlayerQueries.GetPlayerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get_player_by_id_after_update -> %w", err)
	}

	// TODO: End transaction

	return s.getPlayerData(&player)
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

func (s *Player) getOrCreateDailySignInProgress(
	ctx context.Context,
	playerID pgtype.UUID,
) (*playersqlc.MarkDailySignInDaysParams, error) {
	createNewRow := false
	weekStartAt := helper.WeekStartAt(time.Now())

	// TODO: Begin transaction?

	getDailySignInByPlayerIDRow, err := s.PlayerRepo.PlayerQueries.GetDailySignInByPlayerID(
		ctx,
		playersqlc.GetDailySignInByPlayerIDParams{
			PlayerID: playerID,
			WeekStartAt: pgtype.Timestamptz{
				Time:  weekStartAt,
				Valid: true,
			},
		},
	)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get_daily_sign_in_by_player_id -> %w", err)
		}
		createNewRow = true
	}

	var markDailySignInDaysParams playersqlc.MarkDailySignInDaysParams
	if createNewRow {
		initDailySignInRow, err := s.PlayerRepo.PlayerQueries.InitDailySignIn(
			ctx,
			playersqlc.InitDailySignInParams{
				PlayerID: playerID,
				WeekStartAt: pgtype.Timestamptz{
					Time:  weekStartAt,
					Valid: true,
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("init_daily_sign_in -> %w", err)
		}

		markDailySignInDaysParams = playersqlc.MarkDailySignInDaysParams{
			ID:          initDailySignInRow.ID,
			PlayerID:    playerID,
			ClaimedDays: initDailySignInRow.ClaimedDays,
		}
	} else {
		markDailySignInDaysParams = playersqlc.MarkDailySignInDaysParams{
			ID:          getDailySignInByPlayerIDRow.ID,
			PlayerID:    playerID,
			ClaimedDays: getDailySignInByPlayerIDRow.ClaimedDays,
		}
	}

	// TODO: End transaction

	return &markDailySignInDaysParams, nil
}
