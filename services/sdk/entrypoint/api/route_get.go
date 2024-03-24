package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gatekeeper"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseRouteGet
// @Tags			route
// @Router		/route/{id}					[get]
// @Param			id									path			string						true	"endpoint id"
// @Success		200									{object}	RouteGetRes
// @Failure		default							{object}	Error
// @Security	Authorization
func UseRouteGet(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := &usecase.RouteGetIn{
			WsId: r.Context().Value(gatekeeper.CtxTenantId).(string),
			Id:   chi.URLParam(r, "id"),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Route().Get(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &RouteGetRes{Route: &Route{}}
		res.Map(out.Route)
		httpxwriter.Ok(w, res)
	}
}

type RouteGetRes struct {
	*Route
} // @name RouteGetRes
