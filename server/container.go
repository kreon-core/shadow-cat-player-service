package server

import (
	"sc-player-service/controller"
	"sc-player-service/infrastructure/config"
	"sc-player-service/repository"
	"sc-player-service/repository/playersqlc"
	"sc-player-service/service"
	"sc-player-service/temp"
)

type Container struct {
	PlayerSvc   *service.Player
	PlayerRepo  *repository.Player
	PlayerHCtrl *controller.PlayerH
	AuthMW      *temp.Auth
}

func NewContainer(
	cfg *config.Config,
	playerDBQueries *playersqlc.Queries,
) *Container {
	playerRepo := repository.NewPlayer(playerDBQueries)
	playerSvc := service.NewPlayer(playerRepo)
	playerHCtrl := controller.NewPlayerH(playerSvc)
	authMW := temp.NewAuthMiddleware(&cfg.Secrets)

	return &Container{
		PlayerHCtrl: playerHCtrl,
		PlayerSvc:   playerSvc,
		PlayerRepo:  playerRepo,
		AuthMW:      authMW,
	}
}
