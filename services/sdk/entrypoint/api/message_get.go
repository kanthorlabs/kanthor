package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gatekeeper"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseMessageGet
// @Tags			message
// @Router		/message/{id}				[get]
// @Param			id									path			string						true	"message id"
// @Success		200									{object}	MessageGetRes
// @Failure		default							{object}	Error
// @Security	Authorization
func UseMessageGet(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := &usecase.MessageGetIn{
			WsId: r.Context().Value(gatekeeper.CtxTenantId).(string),
			Id:   chi.URLParam(r, "id"),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Message().Get(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &MessageGetRes{Message: &Message{}, Endpoints: make([]*MessageEndpoint, 0)}
		res.Message.Map(out.Message)
		for id := range out.Endpoints {
			msgep := &MessageEndpoint{}
			msgep.Map(out.Endpoints[id])
			res.Endpoints = append(res.Endpoints, msgep)
		}

		httpxwriter.Ok(w, res)
	}
}

type MessageGetRes struct {
	Message   *Message           `json:"message"`
	Endpoints []*MessageEndpoint `json:"endpoints"`
} // @name MessageGetRes
