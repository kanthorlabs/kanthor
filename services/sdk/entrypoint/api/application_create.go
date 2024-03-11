package api

import (
	"encoding/json"
	"net/http"

	"github.com/kanthorlabs/common/gatekeeper"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseApplicationCreate
// @Tags			application
// @Router		/application		[post]
// @Param			request					body			ApplicationCreateReq	true	"request body"
// @Success		200							{object}	ApplicationCreateRes
// @Failure		default					{object}	Error
// @Security	Authorization
func UseApplicationCreate(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ApplicationCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.ErrorString("PORTAl.WORKSPACE.CREATE.DECODE.ERROR"))
			return
		}

		account := r.Context().Value(passport.CtxAccount).(*ppentities.Account)
		in := &usecase.ApplicationCreateIn{
			Modifier: account.Username,
			WsId:     r.Context().Value(gatekeeper.CtxTenantId).(string),
			Name:     req.Name,
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Application().Create(r.Context(), in)
		if err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		res := &ApplicationCreateRes{Application: &Application{}}
		res.Map(out.Application)
		httpxwriter.Ok(w, res)
	}
}

type ApplicationCreateReq struct {
	Name string `json:"name" example:"simple app"`
} // @name ApplicationCreateReq

type ApplicationCreateRes struct {
	*Application
} // @name ApplicationCreateRes
