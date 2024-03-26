package api

import (
	"github.com/go-chi/chi/v5"
	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
)

func RegisterRouteRoutes(router chi.Router, service *sdk) {
	router.Route("/route", func(sr chi.Router) {
		sr.Use(httpxmw.Authz(service.infra.Gatekeeper(), config.ServiceName))
		sr.Post("/", UseRouteCreate(service))
		sr.Get("/", UseRouteList(service))
		sr.Route("/{id}", func(ssr chi.Router) {
			ssr.Get("/", UseRouteGet(service))
			ssr.Patch("/", UseRouteUpdate(service))
			ssr.Delete("/", UseRouteDelete(service))
		})
	})
}

type Route struct {
	Id                  string `json:"id" example:"rt_2dcBT1R8169aIGvx0PEilqrJIYM"`
	CreatedAt           int64  `json:"created_at" example:"1728925200000"`
	UpdatedAt           int64  `json:"updated_at" example:"1728925200000"`
	Name                string `json:"name" example:"passthrough"`
	EpId                string `json:"ep_id" example:"ep_2dZRCcnumVTMI9eHdmep89IpOgY"`
	Priority            int32  `json:"priority" example:"1"`
	Exclusionary        bool   `json:"exclusionary" example:"false"`
	ConditionSource     string `json:"condition_source" example:"type"`
	ConditionExpression string `json:"condition_expression" example:"any::"`
} // @name Route

func (rt *Route) Map(entity *entities.Route) {
	rt.Id = entity.Id
	rt.CreatedAt = entity.CreatedAt
	rt.UpdatedAt = entity.UpdatedAt
	rt.EpId = entity.EpId
	rt.Name = entity.Name
	rt.Priority = entity.Priority
	rt.Exclusionary = entity.Exclusionary
	rt.ConditionSource = entity.ConditionSource
	rt.ConditionExpression = entity.ConditionExpression
}
