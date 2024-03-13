package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseRouteUpdate
// @Tags			route
// @Router		/route/{id}				[patch]
// @Param			ep_id							query			string							true	"endpoint id"
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

		account := r.Context().Value(passport.CtxAccount).(*ppentities.Account)
		ep := r.Context().Value(CtxEndpoint).(*entities.Endpoint)
		in := &usecase.RouteUpdateIn{
			Modifier:            account.Username,
			EpId:                ep.Id,
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
	Name                string `json:"name" example:"only test category route"`
	Priority            int32  `json:"priority" example:"9"`
	Exclusionary        bool   `json:"exclusionary" example:"false"`
	ConditionSource     string `json:"condition_source" example:"category"`
	ConditionExpression string `json:"condition_expression" example:"prefix::test."`
} // @name WorkspaceUpdateReq

type RouteUpdateRes struct {
	*Route
} // @name RouteUpdateRes
