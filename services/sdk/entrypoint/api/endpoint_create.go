package api

import (
	"encoding/json"
	"net/http"

	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseEndpointCreate
// @Tags			endpoint
// @Router		/endpoint			[post]
// @Param			app_id				query			string							true	"application id"
// @Param			request				body			EndpointCreateReq		true	"request body"
// @Success		200						{object}	EndpointCreateRes
// @Failure		default				{object}	Error
// @Security	Authorization
func UseEndpointCreate(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EndpointCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.ErrorString("SDK.ENDPOINT.CREATE.DECODE.ERROR"))
			return
		}

		app := r.Context().Value(CtxApplication).(*entities.Application)
		in := &usecase.EndpointCreateIn{
			AppId:  app.Id,
			Name:   req.Name,
			Method: req.Method,
			Uri:    req.Uri,
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Endpoint().Create(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &EndpointCreateRes{Endpoint: &Endpoint{}}
		res.Map(out.Endpoint)
		httpxwriter.Ok(w, res)
	}
}

type EndpointCreateReq struct {
	Name   string `json:"name" example:"echo endpoint"`
	Method string `json:"method" example:"POST"`
	Uri    string `json:"uri" example:"https://postman-echo.com/post"`
} // @name EndpointCreateReq

type EndpointCreateRes struct {
	*Endpoint
} // @name EndpointCreateRes
