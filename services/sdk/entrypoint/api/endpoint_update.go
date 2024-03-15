package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseEndpointUpdate
// @Tags			endpoint
// @Router		/endpoint/{id}		[patch]
// @Param			app_id						query			string							true	"application id"
// @Param			id								path			string							true	"endpoint id"
// @Param			request						body			EndpointUpdateReq		true	"request body"
// @Success		200								{object}	EndpointUpdateRes
// @Failure		default						{object}	Error
// @Security	Authorization
func UseEndpointUpdate(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EndpointUpdateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.ErrorString("SDK.ENDPOINT.UPDATE.DECODE.ERROR"))
			return
		}

		app := r.Context().Value(CtxApplication).(*entities.Application)
		in := &usecase.EndpointUpdateIn{
			AppId:  app.Id,
			Id:     chi.URLParam(r, "id"),
			Name:   req.Name,
			Method: req.Method,
			Uri:    req.Uri,
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Endpoint().Update(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &EndpointUpdateRes{Endpoint: &Endpoint{}}
		res.Map(out.Endpoint)
		httpxwriter.Ok(w, res)
	}
}

type EndpointUpdateReq struct {
	Name   string `json:"name" example:"echo endpoint with PUT"`
	Method string `json:"method" example:"PUT"`
	Uri    string `json:"uri" example:"https://postman-echo.com/put"`
} // @name WorkspaceUpdateReq

type EndpointUpdateRes struct {
	*Endpoint
} // @name EndpointUpdateRes
