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
	AuthClient *external.HTTPClient

	PlayerRepo *repository.Player

	AuthMW *middleware.Auth

	PlayerSvc *service.Player

	PlayerHCtrl      *controller.PlayerH
	LeaderboardHCtrl *controller.LeaderboardH
	ShopHCtrl        *controller.ShopH
}

func NewContainer(
	cfg *config.Config,
	playerDBQueries *playersqlc.Queries,
) *Container {
	authClient := external.NewClient(&cfg.Externals.AuthClient)

	playerRepo := repository.NewPlayer(playerDBQueries)

	authMW := middleware.NewAuthMiddleware(authClient)

	playerSvc := service.NewPlayer(playerRepo)

	playerHCtrl := controller.NewPlayerH(playerSvc)
	leaderboardHCtrl := controller.NewLeaderboardH()
	shopHCtrl := controller.NewShopH()

	return &Container{
		AuthClient: authClient,

		PlayerRepo: playerRepo,

		AuthMW: authMW,

		PlayerSvc: playerSvc,

		PlayerHCtrl:      playerHCtrl,
		LeaderboardHCtrl: leaderboardHCtrl,
		ShopHCtrl:        shopHCtrl,
	}
}
