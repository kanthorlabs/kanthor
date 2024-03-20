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
			Type:  req.Type,
			Body:  utils.Stringify(req),
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
	Type   string         `json:"type" example:"testing.openapi"`
	Object map[string]any `json:"object" example:"say:hello,from_client:openapi" swaggertype:"object,string"`
} // @name MessageCreateReq

type MessageCreateRes struct {
	Id        string `json:"id" example:"msg_2dgJIHGMePYS4VJRmEGv73RfIvu"`
	CreatedAt int64  `json:"created_at" example:"1728925200000"`
} // @name MessageCreateRes
