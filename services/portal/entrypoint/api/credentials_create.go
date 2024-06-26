package api

import (
	"encoding/json"
	"net/http"

	httpxmw "github.com/kanthorlabs/common/gateway/httpx/middleware"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	passportutils "github.com/kanthorlabs/common/passport/utils"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
	"go.opentelemetry.io/otel"
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
		ctx, span := otel.Tracer("SERVICE.CREDENTIALS.CREATE").Start(r.Context(), "ENTRYPOINT.API")
		defer span.End()

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

		out, err := service.uc.Credentials().Create(ctx, in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &CredentialsCreateRes{
			Tenant:   out.Tenant,
			Username: out.Username,
			Password: out.Password,
			Schemes: httpxwriter.M{
				// the region part will help the client detech what region API it needs to talk to
				"basic": passportutils.CreateRegionalBasicCredentials(out.Username + ":" + out.Password),
			},
		}
		httpxwriter.Ok(w, res)
	}
}

type CredentialsCreateReq struct {
	Name string `json:"name" example:"default credentials"`
} // @name CredentialsCreateReq

type CredentialsCreateRes struct {
	Tenant   string        `json:"tenant" example:"ws_2dXFW6gHgDR9YBPILkfSmnBaCu8"`
	Username string        `json:"username" example:"admin"`
	Password string        `json:"password" example:"b7ccecf6054343ca8c3ebbdc36b05e5bcc28f4b5e812484387ad7de6ad6a04e4"`
	Schemes  httpxwriter.M `json:"schemes" example:"basic:YWRtaW46YjdjY2VmNjA1NDM0M2NhOGMzZWJiZGMzNmIwNWU1YmNjMjhmNGI1ZTgxMjQ4NDM4N2FkN2RlNmFkNmEwNGU0" swaggertype:"object,string"`
} // @name CredentialsCreateRes
