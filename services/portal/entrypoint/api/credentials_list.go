package api

import (
	"net/http"

	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

// UseCredentialsList
// @Tags			credentials
// @Router		/credentials		[get]
// @Success		200							{object}	CredentialsListRes
// @Failure		default					{object}	Error
// @Security	Authorization
// @Security	TenantId
func UseCredentialsList(service *portal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := &usecase.CredentialsListIn{
			Tenant: r.Header.Get(httpxmw.HeaderAuthzTenant),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Credentials().List(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &CredentialsListRes{Data: make([]*CredentialsAccount, len(out.Data))}
		for i := range out.Data {
			res.Data[i] = &CredentialsAccount{}
			res.Data[i].Map(out.Data[i])
		}
		httpxwriter.Ok(w, res)
	}
}

type CredentialsListRes struct {
	Data []*CredentialsAccount `json:"data"`
} // @name CredentialsListRes
