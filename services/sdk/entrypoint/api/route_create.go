package api

import (
	"encoding/json"
	"net/http"

	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseRouteCreate
// @Tags			route
// @Router		/route				[post]
// @Param			ep_id					query			string						true	"endpoint id"
// @Param			request				body			RouteCreateReq		true	"request body"
// @Success		200						{object}	RouteCreateRes
// @Failure		default				{object}	Error
// @Security	Authorization
func UseRouteCreate(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RouteCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.ErrorString("SDK.ENDPOINT.CREATE.DECODE.ERROR"))
			return
		}

		ep := r.Context().Value(CtxEndpoint).(*entities.Endpoint)
		in := &usecase.RouteCreateIn{
			EpId:                ep.Id,
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

		out, err := service.uc.Route().Create(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &RouteCreateRes{Route: &Route{}}
		res.Map(out.Route)
		httpxwriter.Ok(w, res)
	}
}

type RouteCreateReq struct {
	Name                string `json:"name" example:"passthrough"`
	Priority            int32  `json:"priority" example:"1"`
	Exclusionary        bool   `json:"exclusionary" example:"false"`
	ConditionSource     string `json:"condition_source" example:"type"`
	ConditionExpression string `json:"condition_expression" example:"prefix::testing."`
} // @name RouteCreateReq

type RouteCreateRes struct {
	*Route
} // @name RouteCreateRes
