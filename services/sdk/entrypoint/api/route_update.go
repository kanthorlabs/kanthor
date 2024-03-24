package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gatekeeper"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseRouteUpdate
// @Tags			route
// @Router		/route/{id}				[patch]
// @Param			id								path			string							true	"route id"
// @Param			request						body			RouteUpdateReq			true	"request body"
// @Success		200								{object}	RouteUpdateRes
// @Failure		default						{object}	Error
// @Security	Authorization
func UseRouteUpdate(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RouteUpdateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.ErrorString("SDK.ENDPOINT.UPDATE.DECODE.ERROR"))
			return
		}

		in := &usecase.RouteUpdateIn{
			WsId:                r.Context().Value(gatekeeper.CtxTenantId).(string),
			Id:                  chi.URLParam(r, "id"),
			Name:                req.Name,
			Priority:            req.Priority,
			Exclusionary:        req.Exclusionary,
			ConditionSource:     req.ConditionSource,
			ConditionExpression: req.ConditionExpression,
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Route().Update(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &RouteUpdateRes{Route: &Route{}}
		res.Map(out.Route)
		httpxwriter.Ok(w, res)
	}
}

type RouteUpdateReq struct {
	Name                string `json:"name" example:"only test type route"`
	Priority            int32  `json:"priority" example:"9"`
	Exclusionary        bool   `json:"exclusionary" example:"false"`
	ConditionSource     string `json:"condition_source" example:"type"`
	ConditionExpression string `json:"condition_expression" example:"prefix::testing."`
} // @name WorkspaceUpdateReq

type RouteUpdateRes struct {
	*Route
} // @name RouteUpdateRes
