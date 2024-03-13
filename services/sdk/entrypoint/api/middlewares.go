package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gatekeeper"
	"github.com/kanthorlabs/common/gateway"
	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

var (
	CtxApplication gateway.ContextKey = "sdk.application"
	CtxEndpoint    gateway.ContextKey = "sdk.endpoint"
)

// UseApplication help you detect whether the application is belong to the tenant of the request account or not
func UseApplication(service *sdk, param string) httpxmw.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// get from url param first
			id := chi.URLParam(r, param)
			if id == "" {
				// then get from query param
				id = r.URL.Query().Get(param)
			}

			in := &usecase.ApplicationGetIn{
				WsId: ctx.Value(gatekeeper.CtxTenantId).(string),
				Id:   id,
			}
			if err := in.Validate(); err != nil {
				httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
				return
			}

			out, err := service.uc.Application().Get(ctx, in)
			if err != nil {
				httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
				return
			}

			ctx = context.WithValue(ctx, CtxApplication, out.Application)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UseEndpoint help you detect whether the endpoint is belong to the tenant of the request account or not
func UseEndpoint(service *sdk, param string) httpxmw.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// get from url param first
			id := chi.URLParam(r, param)
			if id == "" {
				// then get from query param
				id = r.URL.Query().Get(param)
			}

			in := &usecase.EndpointGetOwnIn{
				WsId: ctx.Value(gatekeeper.CtxTenantId).(string),
				Id:   id,
			}
			if err := in.Validate(); err != nil {
				httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
				return
			}

			out, err := service.uc.Endpoint().GetOwn(ctx, in)
			if err != nil {
				httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
				return
			}

			ctx = context.WithValue(ctx, CtxEndpoint, out.Endpoint)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
