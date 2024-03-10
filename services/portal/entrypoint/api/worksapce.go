package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gateway"
	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

func RegisterWorkspaceRoutes(router chi.Router, service *portal) {
	router.Route("/workspace", func(sr chi.Router) {
		sr.Post("/", UseWorkspaceCreate(service))
		sr.Get("/", UseWorkspaceList(service))
		sr.Route("/{id}", func(ssr chi.Router) {
			ssr.Use(UseWorkspace(service))
			ssr.Use(httpxmw.Authz(service.infra.Gatekeeper(), config.ServiceName))

			ssr.Get("/", UseWorkspaceGet(service))
			ssr.Patch("/", UseWorkspaceUpdate(service))
			ssr.Delete("/", UseWorkspaceDelete(service))
		})
	})
}

var CtxWorksspace gateway.ContextKey = "portal.workspace"

func UseWorkspace(service *portal) httpxmw.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			in := &usecase.WorkspaceGetIn{Id: chi.URLParam(r, "id")}
			out, err := service.uc.Workspace().Get(r.Context(), in)
			if err != nil {
				httpxwriter.ErrNotFound(w, httpxwriter.Error(err))
				return
			}

			// set tenant header for authorization check
			r.Header.Set(httpxmw.HeaderAuthzTenant, out.Id)
			// set workspace context for further use
			r = r.WithContext(context.WithValue(r.Context(), CtxWorksspace, out.Workspace))
			next.ServeHTTP(w, r)
		})
	}
}

type Workspace struct {
	Id        string `json:"id" example:"ws_2nR9p4W6UmUieJMLIf7ilbXBIRR"`
	CreatedAt int64  `json:"created_at" example:"1728925200000"`
	UpdatedAt int64  `json:"updated_at" example:"1728925200000"`
	OwnerId   string `json:"owner_id" example:"admin"`
	Name      string `json:"name" example:"main workspace"`
	Tier      string `json:"tier" example:"default"`
}

func (ws *Workspace) Map(entity *entities.Workspace) {
	ws.Id = entity.Id
	ws.CreatedAt = entity.CreatedAt
	ws.UpdatedAt = entity.UpdatedAt
	ws.OwnerId = entity.OwnerId
	ws.Name = entity.Name
	ws.Tier = entity.Tier
}
