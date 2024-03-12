package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseEndpointGet
// @Tags			endpoint
// @Router		/endpoint/{id}			[get]
// @Param			app_id							query			string						true	"application id"
// @Param			id									path			string						true	"endpoint id"
// @Success		200									{object}	EndpointGetRes
// @Failure		default							{object}	Error
// @Security	Authorization
func UseEndpointGet(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app := r.Context().Value(CtxApplication).(*entities.Application)
		in := &usecase.EndpointGetIn{
			AppId: app.Id,
			Id:    chi.URLParam(r, "id"),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Endpoint().Get(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &EndpointGetRes{Endpoint: &Endpoint{}}
		res.Map(out.Endpoint)
		httpxwriter.Ok(w, res)
	}
}

type EndpointGetRes struct {
	*Endpoint
} // @name EndpointGetRes
