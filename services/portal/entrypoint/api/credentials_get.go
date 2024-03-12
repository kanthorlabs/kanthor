package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

// UseCredentialsGet
// @Tags			credentials
// @Router		/credentials/{username}		[get]
// @Param			username									path			string						true	"credentials username"
// @Success		200												{object}	CredentialsGetRes
// @Failure		default										{object}	Error
// @Security	Authorization
// @Security	TenantId
func UseCredentialsGet(service *portal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := &usecase.CredentialsGetIn{
			Tenant:   r.Header.Get(httpxmw.HeaderAuthzTenant),
			Username: chi.URLParam(r, "username"),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Credentials().Get(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &CredentialsGetRes{&CredentialsAccount{}}
		res.Map(out.CredentialsAccount)
		httpxwriter.Ok(w, res)
	}
}

type CredentialsGetRes struct {
	*CredentialsAccount
} // @name CredentialsGetRes
