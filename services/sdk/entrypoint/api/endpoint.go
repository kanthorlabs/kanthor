package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

func RegisterEndpointRoutes(router chi.Router, service *sdk) {
	router.Route("/endpoint", func(sr chi.Router) {
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
	CreatedAt int64  `json:"created_at" example:"1728925200000"`
	UpdatedAt int64  `json:"updated_at" example:"1728925200000"`
	AppId     string `json:"app_id" example:"app_2dXFXcW6HwrJLQuMjc7n02Xmyq8"`
	Name      string `json:"name" example:"echo endpoint"`
	Method    string `json:"method" example:"POST"`
	Uri       string `json:"uri" example:"https://postman-echo.com/post"`
} // @name Endpoint

func (ep *Endpoint) Map(entity *entities.Endpoint) {
	ep.Id = entity.Id
	ep.CreatedAt = entity.CreatedAt
	ep.UpdatedAt = entity.UpdatedAt
	ep.AppId = entity.AppId
	ep.Name = entity.Name
	ep.Method = entity.Method
	ep.Uri = entity.Uri
}
