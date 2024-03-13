package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseRouteGet
// @Tags			route
// @Router		/route/{id}					[get]
// @Param			ep_id								query			string							true	"endpoint id"
// @Param			id									path			string						true	"endpoint id"
// @Success		200									{object}	RouteGetRes
// @Failure		default							{object}	Error
// @Security	Authorization
func UseRouteGet(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ep := r.Context().Value(CtxEndpoint).(*entities.Endpoint)
		in := &usecase.RouteGetIn{
			EpId: ep.Id,
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
