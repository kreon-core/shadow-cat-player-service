package controller

import (
	"encoding/json"
	"net/http"

	tul "github.com/kreon-core/shadow-cat-common"

	"sc-player-service/helper"
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
	userID, ok := helper.GetFromContext[string](r.Context(), "user_id")
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	data, err := ctrl.PlayerSvc.GetOrCreatePlayer(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get or create player", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := &response.Resp{
		ReturnCode:    tul.Success,
		ReturnMessage: tul.Message(tul.Success),
		Data:          data,
	}
	json.NewEncoder(w).Encode(resp)
}

func (ctrl *PlayerH) Update(w http.ResponseWriter, r *http.Request) {}

func (ctrl *PlayerH) GetEnergy(w http.ResponseWriter, r *http.Request)    {}
func (ctrl *PlayerH) GetInventory(w http.ResponseWriter, r *http.Request) {}

func (ctrl *PlayerH) GetDailyTaskProgress(w http.ResponseWriter, r *http.Request) {}

func (ctrl *PlayerH) GetTowerProgress(w http.ResponseWriter, r *http.Request)   {}
func (ctrl *PlayerH) GetChapterProgress(w http.ResponseWriter, r *http.Request) {}
