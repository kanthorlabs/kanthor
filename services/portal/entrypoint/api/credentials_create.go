package api

import (
	"encoding/json"
	"net/http"

	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

// UseCredentialsCreate
// @Tags			credentials
// @Router		/credentials		[post]
// @Param			request					body			CredentialsCreateReq	true	"request body"
// @Success		200							{object}	CredentialsCreateRes
// @Failure		default					{object}	Error
// @Security	Authorization
// @Security	TenantId
func UseCredentialsCreate(service *portal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CredentialsCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.ErrorString("PORTAl.CREDENTIALS.CREATE.DECODE.ERROR"))
			return
		}

		in := &usecase.CredentialsCreateIn{
			Tenant: r.Header.Get(httpxmw.HeaderAuthzTenant),
			Name:   req.Name,
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Credentials().Create(r.Context(), in)
		if err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		res := &CredentialsCreateRes{
			Tenant:   out.Tenant,
			Username: out.Username,
			Password: out.Password,
		}
		httpxwriter.Ok(w, res)
	}
}

type CredentialsCreateReq struct {
	Name string `json:"name" example:"default credentials"`
} // @name CredentialsCreateReq

type CredentialsCreateRes struct {
	Tenant   string `json:"tenant" example:"ws_2nR9p4W6UmUieJMLIf7ilbXBIRR"`
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"b7ccecf6054343ca8c3ebbdc36b05e5bcc28f4b5e812484387ad7de6ad6a04e4"`
} // @name CredentialsCreateRes