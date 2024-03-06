package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	"github.com/kanthorlabs/kanthor/services/portal/config"
)

func RegisterWorkspaceRoutes(router chi.Router, service *portal) {
	router.Route("/workspace", func(sr chi.Router) {
		sr.Post("/", UseWorkspaceCreate(service))
		sr.Route("/{id}", func(ssr chi.Router) {
			ssr.Use(UseWorkspace())
			ssr.Use(httpxmw.Authz(service.infra.Gatekeeper(), config.ServiceName))
			ssr.Get("/", UseWorkspaceGet(service))
		})
	})
}

func UseWorkspace() httpxmw.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set(httpxmw.HeaderAuthzTenant, chi.URLParam(r, "id"))
			next.ServeHTTP(w, r)
		})
	}
}
