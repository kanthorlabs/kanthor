package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gatekeeper"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseApplicationGet
// @Tags			application
// @Router		/application/{id}		[get]
// @Param			id									path			string						true	"application id"
// @Success		200									{object}	ApplicationGetRes
// @Failure		default							{object}	Error
// @Security	Authorization
func UseApplicationGet(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := &usecase.ApplicationGetIn{
			WsId: r.Context().Value(gatekeeper.CtxTenantId).(string),
			Id:   chi.URLParam(r, "id"),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Application().Get(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &ApplicationGetRes{Application: &Application{}}
		res.Map(out.Application)
		httpxwriter.Ok(w, res)
	}
}

type ApplicationGetRes struct {
	*Application
} // @name ApplicationGetRes
