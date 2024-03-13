package api

import (
	"net/http"

	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	gwquery "github.com/kanthorlabs/common/gateway/query"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseRouteList
// @Tags			route
// @Router		/route						[get]
// @Param			ep_id							query			string						true	"endpoint id"
// @Param			_ids							query			[]string					false	"list by ids"
// @Param			_q								query			string						false	"search keyword"
// @Param			_limit						query			int								false	"limit returning records"	default(5)
// @Param			_page							query			int								false	"current requesting page"	default(0)
// @Success		200								{object}	RouteListRes
// @Failure		default						{object}	Error
// @Security	Authorization
func UseRouteList(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ep := r.Context().Value(CtxEndpoint).(*entities.Endpoint)
		in := &usecase.RouteListIn{
			EpId:  ep.Id,
			Query: gwquery.FromHttpx(r).ToDbPagingQuery(),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Route().List(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &RouteListRes{Count: out.Count, Data: make([]*Route, len(out.Data))}
		for i := range out.Data {
			res.Data[i] = &Route{}
			res.Data[i].Map(out.Data[i])
		}
		httpxwriter.Ok(w, res)
	}
}

type RouteListRes struct {
	Count int64    `json:"count"`
	Data  []*Route `json:"data"`
} // @name RouteListRes
