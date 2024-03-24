package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gatekeeper"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseEndpointDelete
// @Tags			endpoint
// @Router		/endpoint/{id}			[delete]
// @Param			app_id							query			string							true	"application id"
// @Param			id									path			string							true	"endpoint id"
// @Success		200									{object}	EndpointDeleteRes
// @Failure		default							{object}	Error
// @Security	Authorization
func UseEndpointDelete(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := &usecase.EndpointDeleteIn{
			WsId: r.Context().Value(gatekeeper.CtxTenantId).(string),
			Id:   chi.URLParam(r, "id"),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Endpoint().Delete(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &EndpointDeleteRes{Endpoint: &Endpoint{}}
		res.Map(out.Endpoint)
		httpxwriter.Ok(w, res)
	}
}

type EndpointDeleteRes struct {
	*Endpoint
} // @name EndpointDeleteRes
