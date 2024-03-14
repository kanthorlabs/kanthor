package api

import (
	"encoding/json"
	"net/http"

	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

// UseMessageCreate
// @Tags			message
// @Router		/message			[post]
// @Param			app_id				query			string							true	"application id"
// @Param			request				body			MessageCreateReq		true	"request body"
// @Success		200						{object}	MessageCreateRes
// @Failure		default				{object}	Error
// @Security	Authorization
func UseMessageCreate(service *sdk) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MessageCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.ErrorString("SDK.MESSAGE.CREATE.DECODE.ERROR"))
			return
		}

		in := &usecase.MessageCreateIn{
			AppId: r.URL.Query().Get("app_id"),
			Tag:   req.Tag,
			Body:  utils.Stringify(req.Body),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Message().Create(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &MessageCreateRes{Id: out.Id, CreatedAt: out.CreatedAt}
		httpxwriter.Ok(w, res)
	}
}

type MessageCreateReq struct {
	Tag  string         `json:"tag" example:"testing.openapi"`
	Body map[string]any `json:"uri"  swaggertype:"object,string" example:"say:hello,from:openapi"`
} // @name MessageCreateReq

type MessageCreateRes struct {
	Id        string `json:"id" example:"msg_2dgJIHGMePYS4VJRmEGv73RfIvu"`
	CreatedAt int64  `json:"created_at" example:"1728925200000"`
} // @name MessageCreateRes
