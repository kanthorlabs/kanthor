package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

// UseCredentialsUpdate
// @Tags			credentials
// @Router		/credentials/{username}		[patch]
// @Param			username									path			string						true	"credentials username"
// @Param			request										body			CredentialsUpdateReq	true	"request body"
// @Success		200												{object}	CredentialsUpdateRes
// @Failure		default										{object}	Error
// @Security	Authorization
// @Security	TenantId
func UseCredentialsUpdate(service *portal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CredentialsUpdateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.ErrorString("PORTAl.CREDENTIALS.UPDATE.DECODE.ERROR"))
			return
		}

		in := &usecase.CredentialsUpdateIn{
			Tenant:   r.Header.Get(httpxmw.HeaderAuthzTenant),
			Username: chi.URLParam(r, "username"),
			Name:     req.Name,
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Credentials().Update(r.Context(), in)
		if err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		res := &CredentialsUpdateRes{&CredentialsAccount{}}
		res.Map(out.CredentialsAccount)
		httpxwriter.Ok(w, res)
	}
}

type CredentialsUpdateReq struct {
	Name string `json:"name" example:"another name"`
} // @name CredentialsUpdateReq

type CredentialsUpdateRes struct {
	*CredentialsAccount
} // @name CredentialsUpdateRes
