package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

func RegisterApplicationRoutes(router chi.Router, service *sdk) {
	router.Route("/application", func(sr chi.Router) {
		sr.Post("/", UseApplicationCreate(service))
		sr.Get("/", UseApplicationList(service))
		sr.Route("/{id}", func(ssr chi.Router) {
			ssr.Get("/", UseApplicationGet(service))
			ssr.Patch("/", UseApplicationUpdate(service))
			ssr.Delete("/", UseApplicationDelete(service))
		})
	})
}

type Application struct {
	Id        string `json:"id" example:"app_2dXFXcW6HwrJLQuMjc7n02Xmyq8"`
	WsId      string `json:"ws_id" example:"ws_2dXFW6gHgDR9YBPILkfSmnBaCu8"`
	CreatedAt int64  `json:"created_at" example:"1728925200000"`
	UpdatedAt int64  `json:"updated_at" example:"1728925200000"`
	Name      string `json:"name" example:"main application"`
} // @name Application

func (ws *Application) Map(entity *entities.Application) {
	ws.Id = entity.Id
	ws.WsId = entity.WsId
	ws.CreatedAt = entity.CreatedAt
	ws.UpdatedAt = entity.UpdatedAt
	ws.Name = entity.Name
}
