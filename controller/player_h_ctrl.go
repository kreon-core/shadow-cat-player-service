package controller

import (
	"encoding/json"
	"net/http"

	"github.com/kreon-core/shadow-cat-common/appc"
	"github.com/kreon-core/shadow-cat-common/ctxc"
	"github.com/kreon-core/shadow-cat-common/logc"
	"github.com/kreon-core/shadow-cat-common/resc"

	"sc-player-service/helper"
	"sc-player-service/model/api/request"
	"sc-player-service/model/api/response"
	"sc-player-service/service"
)

type PlayerH struct {
	PlayerSvc *service.Player
}

func NewPlayerH(playerSvc *service.Player) *PlayerH {
	return &PlayerH{
		PlayerSvc: playerSvc,
	}
}

func (ctrl *PlayerH) Get(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_Get: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	data, err := ctrl.PlayerSvc.GetOrCreatePlayer(r.Context(), playerID)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_Get: Failed to get or create player")
		unspecifiedErrorResponse(w)
		return
	}

	err = ctrl.PlayerSvc.MarkDailySignIn(r.Context(), playerID)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_Get: Failed to mark daily sign-in")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) Update(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_Update: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	var req request.UpdatePlayer
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_Update: Failed to decode request body")
		badRequestResponse(w)
		return
	}
	defer func() { _ = r.Body.Close() }()

	data, err := ctrl.PlayerSvc.UpdatePlayer(r.Context(), playerID, &req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_Update: Failed to update player")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) GetEnergy(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_Update: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	data, err := ctrl.PlayerSvc.GetEnergy(r.Context(), playerID)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_GetEnergy: Failed to get player energy")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) GetInventory(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_GetInventory: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	data, err := ctrl.PlayerSvc.GetInventory(r.Context(), playerID)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_GetInventory: Failed to get player inventory")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) GetTowerProgress(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_GetTowerProgress: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	data, err := ctrl.PlayerSvc.GetTowerProgress(r.Context(), playerID)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_GetTowerProgress: Failed to get player tower progress")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) GetChapterProgress(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_GetChapterProgress: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	data, err := ctrl.PlayerSvc.GetChapterProgress(r.Context(), playerID)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_GetChapterProgress: Failed to get player chapter progress")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) ClaimChapterRewards(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_ClaimChapterRewards: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	var req request.ClaimChapterRewards
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_ClaimChapterRewards: Failed to decode request body")
		badRequestResponse(w)
		return
	}
	defer func() { _ = r.Body.Close() }()

	data, err := ctrl.PlayerSvc.ClaimChapterRewards(r.Context(), playerID, &req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_ClaimChapterRewards: Failed to claim chapter rewards")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) GetDailySignInProgress(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_GetDailySignInProgress: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	data, err := ctrl.PlayerSvc.GetDailySignInProgress(r.Context(), playerID)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_GetDailySignInProgress: Failed to get player daily sign-in progress")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) UnlockDailySignIn(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_UnlockDailySignIn: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	var req request.UnlockDailySignIn
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_UnlockDailySignIn: Failed to decode request body")
		badRequestResponse(w)
		return
	}
	defer func() { _ = r.Body.Close() }()

	data, err := ctrl.PlayerSvc.UnlockDailySignIn(r.Context(), playerID, &req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_UnlockDailySignIn: Failed to unlock daily sign-in")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) ClaimDailySignRewards(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_ClaimDailySignRewards: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	var req request.ClaimDailySignIn
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_ClaimDailySignRewards: Failed to decode request body")
		badRequestResponse(w)
		return
	}
	defer func() { _ = r.Body.Close() }()

	data, err := ctrl.PlayerSvc.ClaimDailySignIn(r.Context(), playerID, &req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_ClaimDailySignRewards: Failed to claim daily sign-in rewards")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) GetDailyTaskProgress(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_GetDailyTaskProgress: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	data, err := ctrl.PlayerSvc.GetDailyTaskProgress(r.Context(), playerID)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_GetDailyTaskProgress: Failed to get player daily task progress")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) ClaimDailyTaskRewards(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_ClaimDailyTaskRewards: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	var req request.ClaimDailyTask
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_ClaimDailyTaskRewards: Failed to decode request body")
		badRequestResponse(w)
		return
	}
	defer func() { _ = r.Body.Close() }()

	data, err := ctrl.PlayerSvc.ClaimDailyTask(r.Context(), playerID, &req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_ClaimDailyTaskRewards: Failed to claim daily task rewards")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) ExchangeGemsForCoins(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_ExchangeGemsForCoins: Unable to get player ID from context")
		unauthorizedResponse(w)
		return
	}

	var req request.ExchangeGemsForCoins
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_ExchangeGemsForCoins: Failed to decode request body")
		badRequestResponse(w)
		return
	}
	defer func() { _ = r.Body.Close() }()

	data, err := ctrl.PlayerSvc.ExchangeGemsForCoins(r.Context(), playerID, &req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_ExchangeGemsForCoins: Failed to exchange gems for coins")
		unspecifiedErrorResponse(w)
		return
	}

	successResponse(w, data)
}

func (ctrl *PlayerH) GetDailyShopProgress(w http.ResponseWriter, r *http.Request)   {}
func (ctrl *PlayerH) PurchaseDailyShopItems(w http.ResponseWriter, r *http.Request) {}

func (ctrl *PlayerH) StartBattle(w http.ResponseWriter, r *http.Request)   {}
func (ctrl *PlayerH) ResumeBattle(w http.ResponseWriter, r *http.Request)  {}
func (ctrl *PlayerH) BuyBattleLife(w http.ResponseWriter, r *http.Request) {}
func (ctrl *PlayerH) ExitBattle(w http.ResponseWriter, r *http.Request)    {}

func unauthorizedResponse(w http.ResponseWriter) {
	resc.JSON(w, http.StatusUnauthorized, &response.Resp{
		ReturnCode:    appc.EInvalidAccessToken,
		ReturnMessage: appc.Message(appc.EInvalidAccessToken),
	})
}

func badRequestResponse(w http.ResponseWriter) {
	resc.JSON(w, http.StatusBadRequest, &response.Resp{
		ReturnCode:    appc.EInvalidRequest,
		ReturnMessage: appc.Message(appc.EInvalidRequest),
	})
}

func unspecifiedErrorResponse(w http.ResponseWriter) {
	resc.JSON(w, http.StatusBadRequest, &response.Resp{
		ReturnCode:    appc.UUnspecifiedError,
		ReturnMessage: appc.Message(appc.UUnspecifiedError),
	})
}

func successResponse(w http.ResponseWriter, data any) {
	resc.JSON(w, http.StatusOK, &response.Resp{
		ReturnCode:    appc.Success,
		ReturnMessage: appc.Message(appc.Success),
		Data:          data,
	})
}
