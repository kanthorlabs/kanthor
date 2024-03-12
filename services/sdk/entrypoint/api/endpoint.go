package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

func RegisterEndpointRoutes(router chi.Router, service *sdk) {
	router.Route("/endpoint", func(sr chi.Router) {
		sr.Use(UseApplication(service, "app_id"))
		sr.Post("/", UseEndpointCreate(service))
		sr.Get("/", UseEndpointList(service))
		sr.Route("/{id}", func(ssr chi.Router) {
			ssr.Get("/", UseEndpointGet(service))
			ssr.Patch("/", UseEndpointUpdate(service))
			ssr.Delete("/", UseEndpointDelete(service))
		})
	})
}

type Endpoint struct {
	Id        string `json:"id" example:"ep_2dZRCcnumVTMI9eHdmep89IpOgY"`
	AppId     string `json:"app_id" example:"app_2dXFXcW6HwrJLQuMjc7n02Xmyq8"`
	CreatedAt int64  `json:"created_at" example:"1728925200000"`
	UpdatedAt int64  `json:"updated_at" example:"1728925200000"`
	Name      string `json:"name" example:"echo endpoint"`
	Method    string `json:"method" example:"POST"`
	Uri       string `json:"uri" example:"https://postman-echo.com/post"`
} // @name Endpoint

func (ws *Endpoint) Map(entity *entities.Endpoint) {
	ws.Id = entity.Id
	ws.AppId = entity.AppId
	ws.CreatedAt = entity.CreatedAt
	ws.UpdatedAt = entity.UpdatedAt
	ws.Name = entity.Name
	ws.Method = entity.Method
	ws.Uri = entity.Uri
}
