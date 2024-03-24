package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gatekeeper"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseEndpointGetSecret
// @Tags			endpoint
// @Router		/endpoint/{id}/secret			[get]
// @Param			id												path			string						true	"endpoint id"
// @Success		200												{object}	EndpointGetSecretRes
// @Failure		default										{object}	Error
// @Security	Authorization
func UseEndpointGetSecretSecret(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := &usecase.EndpointGetSecretIn{
			WsId: r.Context().Value(gatekeeper.CtxTenantId).(string),
			Id:   chi.URLParam(r, "id"),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Endpoint().GetSecret(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &EndpointGetSecretRes{Endpoint: &Endpoint{}, SecretKey: out.SecretKey}
		res.Map(out.Endpoint)
		httpxwriter.Ok(w, res)
	}
}

type EndpointGetSecretRes struct {
	*Endpoint
	SecretKey string
} // @name EndpointGetSecretRes
