package server

import (
	"sc-player-service/controller"
	"sc-player-service/repository"
	"sc-player-service/repository/pgsqlc"
	"sc-player-service/service"
)

type Container struct {
	PlayerSvc *service.Player

	PlayerRepo *repository.Player

	PlayerHCtrl *controller.PlayerH
}

func NewContainer(
	playerDBQueries *pgsqlc.Queries,
) *Container {
	playerRepo := repository.NewPlayer(playerDBQueries)

	playerSvc := service.NewPlayer(playerRepo)

	playerHCtrl := controller.NewPlayerH(playerSvc)

	return &Container{
		PlayerHCtrl: playerHCtrl,
		PlayerSvc:   playerSvc,
		PlayerRepo:  playerRepo,
	}
}
