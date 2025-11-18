package controller

import "net/http"

type LeaderboardH struct{}

func NewLeaderboardH() *LeaderboardH {
	return &LeaderboardH{}
}

func (ctrl *LeaderboardH) GetMapLeaderboard(w http.ResponseWriter, r *http.Request) {}

func (ctrl *LeaderboardH) GetTowerLeaderboard(w http.ResponseWriter, r *http.Request) {}
