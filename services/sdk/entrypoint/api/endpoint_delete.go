package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
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
		account := r.Context().Value(passport.CtxAccount).(*ppentities.Account)
		app := r.Context().Value(CtxApplication).(*entities.Application)
		in := &usecase.EndpointDeleteIn{
			Modifier: account.Username,
			AppId:    app.Id,
			Id:       chi.URLParam(r, "id"),
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