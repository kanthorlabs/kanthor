package api

import (
	"github.com/go-chi/chi/v5"
)

func RegisterWorkspaceRoutes(router chi.Router, service *portal) {
	router.Route("/workspace", func(sr chi.Router) {
		sr.Post("/", UseWorkspaceCreate(service))
	})
}
