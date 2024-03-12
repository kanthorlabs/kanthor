package api

import (
	"github.com/go-chi/chi/v5"
	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/portal/config"
)

func RegisterWorkspaceRoutes(router chi.Router, service *portal) {
	router.Route("/workspace", func(sr chi.Router) {
		sr.Post("/", UseWorkspaceCreate(service))
		sr.Get("/", UseWorkspaceList(service))
		sr.Route("/{id}", func(ssr chi.Router) {
			ssr.Use(UseWorkspace(service, "id"))
			ssr.Use(httpxmw.Authz(service.infra.Gatekeeper(), config.ServiceName))

			ssr.Get("/", UseWorkspaceGet(service))
			ssr.Patch("/", UseWorkspaceUpdate(service))
			ssr.Delete("/", UseWorkspaceDelete(service))
		})
	})
}

type Workspace struct {
	Id        string `json:"id" example:"ws_2dXFW6gHgDR9YBPILkfSmnBaCu8"`
	CreatedAt int64  `json:"created_at" example:"1728925200000"`
	UpdatedAt int64  `json:"updated_at" example:"1728925200000"`
	OwnerId   string `json:"owner_id" example:"admin"`
	Name      string `json:"name" example:"main workspace"`
	Tier      string `json:"tier" example:"default"`
} // @name Workspace

func (ws *Workspace) Map(entity *entities.Workspace) {
	ws.Id = entity.Id
	ws.CreatedAt = entity.CreatedAt
	ws.UpdatedAt = entity.UpdatedAt
	ws.OwnerId = entity.OwnerId
	ws.Name = entity.Name
	ws.Tier = entity.Tier
}
