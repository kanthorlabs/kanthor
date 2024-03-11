package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

// UseCredentialsExpire
// @Tags			credentials
// @Router		/credentials/{username}/expiration		[put]
// @Param			username															path			string						true	"credentials username"
// @Param			request																body			CredentialsExpireReq	true	"request body"
// @Success		200																		{object}	CredentialsExpireRes
// @Failure		default																{object}	Error
// @Security	Authorization
// @Security	TenantId
func UseCredentialsExpire(service *portal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CredentialsExpireReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.ErrorString("PORTAl.CREDENTIALS.EXPIRE.DECODE.ERROR"))
			return
		}

		in := &usecase.CredentialsExpireIn{
			Tenant:    r.Header.Get(httpxmw.HeaderAuthzTenant),
			Username:  chi.URLParam(r, "username"),
			ExpiresIn: req.ExpiresIn,
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Credentials().Expire(r.Context(), in)
		if err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		res := &CredentialsExpireRes{&CredentialsAccount{}}
		res.Map(out.CredentialsAccount)
		httpxwriter.Ok(w, res)
	}
}

type CredentialsExpireReq struct {
	ExpiresIn int64 `json:"expires_in" example:"1800000"`
} // @name CredentialsCreateReq

type CredentialsExpireRes struct {
	*CredentialsAccount
} // @name CredentialsExpireRes
