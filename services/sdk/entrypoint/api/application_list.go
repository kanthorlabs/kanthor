package api

import (
	"net/http"

	"github.com/kanthorlabs/common/gatekeeper"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	gwquery "github.com/kanthorlabs/common/gateway/query"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseApplicationList
// @Tags			application
// @Router		/application			[get]
// @Param			_ids							query			[]string					false	"list by ids"
// @Param			_q								query			string						false	"search keyword"
// @Param			_limit						query			int								false	"limit returning records"	default(5)
// @Param			_page							query			int								false	"current requesting page"	default(0)
// @Success		200								{object}	ApplicationListRes
// @Failure		default						{object}	Error
// @Security	Authorization
func UseApplicationList(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := &usecase.ApplicationListIn{
			WsId:  r.Context().Value(gatekeeper.CtxTenantId).(string),
			Query: gwquery.FromHttpx(r).ToDbPagingQuery(),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Application().List(r.Context(), in)
		if err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		res := &ApplicationListRes{Count: out.Count, Data: make([]*Application, len(out.Data))}
		for i := range out.Data {
			res.Data[i] = &Application{}
			res.Data[i].Map(out.Data[i])
		}
		httpxwriter.Ok(w, res)
	}
}

type ApplicationListRes struct {
	Count int64          `json:"count"`
	Data  []*Application `json:"data"`
} // @name ApplicationListRes
