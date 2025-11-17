package controller

import (
	"net/http"

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

func (ctrl *PlayerH) Get(w http.ResponseWriter, r *http.Request)    {}
func (ctrl *PlayerH) Update(w http.ResponseWriter, r *http.Request) {}

func (ctrl *PlayerH) GetEnergy(w http.ResponseWriter, r *http.Request)    {}
func (ctrl *PlayerH) GetInventory(w http.ResponseWriter, r *http.Request) {}

func (ctrl *PlayerH) GetDailyTaskProgress(w http.ResponseWriter, r *http.Request) {}

func (ctrl *PlayerH) GetTowerProgress(w http.ResponseWriter, r *http.Request)   {}
func (ctrl *PlayerH) GetChapterProgress(w http.ResponseWriter, r *http.Request) {}
