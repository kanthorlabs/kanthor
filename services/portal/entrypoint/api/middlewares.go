package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gateway"
	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
)

var CtxWorksspace gateway.ContextKey = "portal.workspace"

func UseWorkspace(service *portal, param string) httpxmw.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get from url param first
			id := chi.URLParam(r, param)
			if id == "" {
				// then get from query param
				id = r.URL.Query().Get(param)
			}

			// set tenant header for authorization check
			r.Header.Set(httpxmw.HeaderAuthzTenant, id)
			next.ServeHTTP(w, r)
		})
	}
}
