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
		resc.JSON(w, http.StatusUnauthorized, &response.Resp{
			ReturnCode:    appc.EInvalidAccessToken,
			ReturnMessage: appc.Message(appc.EInvalidAccessToken),
		})
		return
	}

	data, err := ctrl.PlayerSvc.GetOrCreatePlayer(r.Context(), playerID)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_Get: Failed to get or create player")
		resc.JSON(w, http.StatusInternalServerError, &response.Resp{
			ReturnCode:    appc.UUnspecifiedError,
			ReturnMessage: appc.Message(appc.UUnspecifiedError),
		})
		return
	}

	resc.JSON(w, http.StatusOK, &response.Resp{
		ReturnCode:    appc.Success,
		ReturnMessage: appc.Message(appc.Success),
		Data:          data,
	})
}

func (ctrl *PlayerH) Update(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_Update: Unable to get player ID from context")
		resc.JSON(w, http.StatusUnauthorized, &response.Resp{
			ReturnCode:    appc.EInvalidAccessToken,
			ReturnMessage: appc.Message(appc.EInvalidAccessToken),
		})
		return
	}

	var req request.UpdatePlayer
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_Update: Failed to decode request body")
		resc.JSON(w, http.StatusBadRequest, &response.Resp{
			ReturnCode:    appc.EInvalidRequest,
			ReturnMessage: appc.Message(appc.EInvalidRequest),
		})
		return
	}
	defer func() { _ = r.Body.Close() }()

	data, err := ctrl.PlayerSvc.UpdatePlayer(r.Context(), playerID, &req)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_Update: Failed to update player")
		resc.JSON(w, http.StatusInternalServerError, &response.Resp{
			ReturnCode:    appc.UUnspecifiedError,
			ReturnMessage: appc.Message(appc.UUnspecifiedError),
		})
		return
	}

	resc.JSON(w, http.StatusOK, &response.Resp{
		ReturnCode:    appc.Success,
		ReturnMessage: appc.Message(appc.Success),
		Data:          data,
	})
}

func (ctrl *PlayerH) GetEnergy(w http.ResponseWriter, r *http.Request) {
	playerID, ok := ctxc.GetFromContext[string](r.Context(), helper.PlayerIDContextKey)
	if !ok {
		logc.Error().Msg("PlayerH_Update: Unable to get player ID from context")
		resc.JSON(w, http.StatusUnauthorized, &response.Resp{
			ReturnCode:    appc.EInvalidAccessToken,
			ReturnMessage: appc.Message(appc.EInvalidAccessToken),
		})
		return
	}

	data, err := ctrl.PlayerSvc.GetEnergy(r.Context(), playerID)
	if err != nil {
		logc.Error().Err(err).Msg("PlayerH_GetEnergy: Failed to get player energy")
		resc.JSON(w, http.StatusInternalServerError, &response.Resp{
			ReturnCode:    appc.UUnspecifiedError,
			ReturnMessage: appc.Message(appc.UUnspecifiedError),
		})
		return
	}

	resc.JSON(w, http.StatusOK, &response.Resp{
		ReturnCode:    appc.Success,
		ReturnMessage: appc.Message(appc.Success),
		Data:          data,
	})
}

func (ctrl *PlayerH) GetInventory(w http.ResponseWriter, r *http.Request) {}

func (ctrl *PlayerH) GetDailyTaskProgress(w http.ResponseWriter, r *http.Request) {}

func (ctrl *PlayerH) GetTowerProgress(w http.ResponseWriter, r *http.Request)   {}
func (ctrl *PlayerH) GetChapterProgress(w http.ResponseWriter, r *http.Request) {}
