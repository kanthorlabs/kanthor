package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseEndpointGetSecret
// @Tags			endpoint
// @Router		/endpoint/{id}/secret			[get]
// @Param			app_id										query			string						true	"application id"
// @Param			id												path			string						true	"endpoint id"
// @Success		200												{object}	EndpointGetSecretRes
// @Failure		default										{object}	Error
// @Security	Authorization
func UseEndpointGetSecretSecret(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app := r.Context().Value(CtxApplication).(*entities.Application)
		in := &usecase.EndpointGetSecretIn{
			AppId: app.Id,
			Id:    chi.URLParam(r, "id"),
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
