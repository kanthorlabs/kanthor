package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gatekeeper"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseApplicationDelete
// @Tags			application
// @Router		/application/{id}		[delete]
// @Param			id									path			string						true	"application id"
// @Success		200									{object}	ApplicationDeleteRes
// @Failure		default							{object}	Error
// @Security	Authorization
func UseApplicationDelete(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account := r.Context().Value(passport.CtxAccount).(*ppentities.Account)
		in := &usecase.ApplicationDeleteIn{
			Modifier: account.Username,
			WsId:     r.Context().Value(gatekeeper.CtxTenantId).(string),
			Id:       chi.URLParam(r, "id"),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Application().Delete(r.Context(), in)
		if err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		res := &ApplicationDeleteRes{Application: &Application{}}
		res.Map(out.Application)
		httpxwriter.Ok(w, res)
	}
}

type ApplicationDeleteRes struct {
	*Application
} // @name ApplicationDeleteRes
