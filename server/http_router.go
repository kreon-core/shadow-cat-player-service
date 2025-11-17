package server

import "github.com/go-chi/chi/v5"

func LoadRoutes(container *Container) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/players", func(r chi.Router) {
		})
	}
}
