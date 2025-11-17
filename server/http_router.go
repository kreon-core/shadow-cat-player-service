package server

import "github.com/go-chi/chi/v5"

func LoadRoutes(container *Container) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(container.AuthMW.Handle)
		r.Route("/player", func(r chi.Router) {
			r.Get("/", container.PlayerHCtrl.Get)
			r.Put("/", container.PlayerHCtrl.Update)

			r.Get("/energy", container.PlayerHCtrl.GetEnergy)
			r.Get("/inventory", container.PlayerHCtrl.GetInventory)

			r.Get("/tower-progress", container.PlayerHCtrl.GetTowerProgress)
			r.Get("/daily-task-progress", container.PlayerHCtrl.GetDailyTaskProgress)
			r.Get("/chapter-progress", container.PlayerHCtrl.GetChapterProgress)
		})
	}
}
