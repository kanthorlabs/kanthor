package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gatekeeper"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseApplicationUpdate
// @Tags			application
// @Router		/application/{id}		[patch]
// @Param			id									path			string								true	"application id"
// @Param			request							body			ApplicationUpdateReq	true	"request body"
// @Success		200									{object}	ApplicationUpdateRes
// @Failure		default							{object}	Error
// @Security	Authorization
func UseApplicationUpdate(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ApplicationUpdateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.ErrorString("SDK.APPLICATION.UPDATE.DECODE.ERROR"))
			return
		}

		account := r.Context().Value(passport.CtxAccount).(*ppentities.Account)
		in := &usecase.ApplicationUpdateIn{
			Modifier: account.Username,
			WsId:     r.Context().Value(gatekeeper.CtxTenantId).(string),
			Id:       chi.URLParam(r, "id"),
			Name:     req.Name,
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Application().Update(r.Context(), in)
		if err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		res := &ApplicationUpdateRes{Application: &Application{}}
		res.Map(out.Application)
		httpxwriter.Ok(w, res)
	}
}

type ApplicationUpdateReq struct {
	Name string `json:"name" example:"anthor application name"`
} // @name WorkspaceUpdateReq

type ApplicationUpdateRes struct {
	*Application
} // @name ApplicationUpdateRes
