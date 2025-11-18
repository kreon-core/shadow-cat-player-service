package server

import (
	"sc-player-service/controller"
	"sc-player-service/infrastructure/config"
	"sc-player-service/infrastructure/external"
	"sc-player-service/middleware"
	"sc-player-service/repository"
	"sc-player-service/repository/playersqlc"
	"sc-player-service/service"
)

type Container struct {
	AuthClient  *external.HTTPClient
	AuthMW      *middleware.Auth
	PlayerSvc   *service.Player
	PlayerRepo  *repository.Player
	PlayerHCtrl *controller.PlayerH
}

func NewContainer(
	cfg *config.Config,
	playerDBQueries *playersqlc.Queries,
) *Container {
	authClient := external.NewClient(&cfg.Externals.AuthClient)
	authMW := middleware.NewAuthMiddleware(authClient)
	playerRepo := repository.NewPlayer(playerDBQueries)
	playerSvc := service.NewPlayer(playerRepo)
	playerHCtrl := controller.NewPlayerH(playerSvc)

	return &Container{
		AuthClient:  authClient,
		AuthMW:      authMW,
		PlayerHCtrl: playerHCtrl,
		PlayerSvc:   playerSvc,
		PlayerRepo:  playerRepo,
	}
}
