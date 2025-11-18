package server

import (
	"sc-player-service/controller"
	"sc-player-service/infrastructure/config"
	"sc-player-service/repository"
	"sc-player-service/repository/playersqlc"
	"sc-player-service/service"
)

type Container struct {
	PlayerSvc   *service.Player
	PlayerRepo  *repository.Player
	PlayerHCtrl *controller.PlayerH
}

func NewContainer(
	cfg *config.Config,
	playerDBQueries *playersqlc.Queries,
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
