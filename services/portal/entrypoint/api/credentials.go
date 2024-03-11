package api

import (
	"github.com/go-chi/chi/v5"
	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/kanthor/services/permissions"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

// RegisterCredentialsRoutes registers the credentials routes that is sub-routed under the workspace router
func RegisterCredentialsRoutes(router chi.Router, service *portal) {
	router.Route("/credentials", func(sr chi.Router) {
		sr.Use(httpxmw.Authz(service.infra.Gatekeeper(), permissions.Owner))
		sr.Post("/", UseCredentialsCreate(service))
		sr.Get("/", UseCredentialsList(service))
		sr.Route("/{username}", func(ssr chi.Router) {
			ssr.Get("/", UseCredentialsGet(service))
			ssr.Patch("/", UseCredentialsUpdate(service))
			ssr.Put("/expiration", UseCredentialsExpire(service))
		})
	})
}

type CredentialsAccount struct {
	Username      string         `json:"username"`
	Roles         []string       `json:"role"`
	Name          string         `json:"name"`
	Metadata      *safe.Metadata `json:"metadata"`
	CreatedAt     int64          `json:"created_at"`
	UpdatedAt     int64          `json:"updated_at"`
	DeactivatedAt int64          `json:"deactivated_at"`
} // @name CredentialsAccount

func (ws *CredentialsAccount) Map(entity *usecase.CredentialsAccount) {
	ws.Username = entity.Username
	ws.Roles = entity.Roles
	ws.Name = entity.Name
	ws.Metadata = entity.Metadata
	ws.CreatedAt = entity.CreatedAt
	ws.UpdatedAt = entity.UpdatedAt
	ws.DeactivatedAt = entity.DeactivatedAt
}
