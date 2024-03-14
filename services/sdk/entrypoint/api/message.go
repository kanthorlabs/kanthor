package api

import (
	"github.com/go-chi/chi/v5"
)

func RegisterMessageRoutes(router chi.Router, service *sdk) {
	router.Route("/message", func(sr chi.Router) {
		// this API need achieve the best performance,
		// so we pass the application verification into the handler
		// By that way, we can apply cache technique to the application verification
		sr.Post("/", UseMessageCreate(service))
	})
}
