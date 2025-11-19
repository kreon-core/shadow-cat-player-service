package server

import (
	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"
)

func LoadRoutes(container *Container) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(chiMW.AllowContentType("application/json"))
		r.Use(container.AuthMW.VerifyUser)
		r.Route("/player", func(r chi.Router) {
			r.Get("/", container.PlayerHCtrl.Get)
			r.Put("/", container.PlayerHCtrl.Update)

			r.Get("/energy", container.PlayerHCtrl.GetEnergy)
			r.Post("/consume-energy", container.PlayerHCtrl.ConsumeEnergy)
			r.Post("/recharge-energy", container.PlayerHCtrl.RechargeEnergy)

			r.Get("/inventory", container.PlayerHCtrl.GetInventory)
			r.Post("/unlock-new-skins", container.PlayerHCtrl.UnlockNewSkins)
			r.Post("/gain-new-props", container.PlayerHCtrl.GainNewProps)

			r.Get("/tower-progress", container.PlayerHCtrl.GetTowerProgress)
			r.Put("/tower-progress", container.PlayerHCtrl.UpdateTowerProgress)

			r.Get("/chapter-progress", container.PlayerHCtrl.GetChapterProgress)
			r.Post("/claim-chapter-checkpoint", container.PlayerHCtrl.ClaimChapterRewards)

			r.Get("/daily-sign-in-progress", container.PlayerHCtrl.GetDailySignInProgress)
			r.Post("/unlock-daily-sign-in-bonus", container.PlayerHCtrl.UnlockDailySignIn)
			r.Post("/claim-daily-sign-in-rewards", container.PlayerHCtrl.ClaimDailySignInRewards)

			r.Get("/daily-tasks-progress", container.PlayerHCtrl.GetDailyTaskProgress)
			r.Put("/daily-tasks-progress", container.PlayerHCtrl.UpdateDailyTaskProgress)
			r.Post("/claim-daily-tasks-rewards", container.PlayerHCtrl.ClaimDailyTaskRewards)

			r.Post("/exchange-currency", container.PlayerHCtrl.ExchangeGemsForCoins)

			r.Post("/start-battle", container.PlayerHCtrl.StartBattle)
			r.Post("/resume-battle", container.PlayerHCtrl.ResumeBattle)
			r.Post("/complete-battle", container.PlayerHCtrl.CompleteBattle)
			r.Post("/exit-battle", container.PlayerHCtrl.ExitBattle)
		})
	}
}
