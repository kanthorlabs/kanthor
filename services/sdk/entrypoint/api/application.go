package api

import (
	"github.com/go-chi/chi/v5"
	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
)

func RegisterApplicationRoutes(router chi.Router, service *sdk) {
	router.Route("/application", func(sr chi.Router) {
		sr.Use(httpxmw.Authz(service.infra.Gatekeeper(), config.ServiceName))

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
	CreatedAt int64  `json:"created_at" example:"1728925200000"`
	UpdatedAt int64  `json:"updated_at" example:"1728925200000"`
	WsId      string `json:"ws_id" example:"ws_2dXFW6gHgDR9YBPILkfSmnBaCu8"`
	Name      string `json:"name" example:"main application"`
} // @name Application

func (app *Application) Map(entity *entities.Application) {
	app.Id = entity.Id
	app.CreatedAt = entity.CreatedAt
	app.UpdatedAt = entity.UpdatedAt
	app.WsId = entity.WsId
	app.Name = entity.Name
}
