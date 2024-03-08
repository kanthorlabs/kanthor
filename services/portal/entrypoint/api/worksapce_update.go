package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

// UseWorkspaceUpdate
// @Tags			workspace
// @Router		/workspace/{id}		[patch]
// @Param			id								path			string							true	"workspace id"
// @Param			request						body			WorkspaceUpdateReq	true	"request body"
// @Success		200								{object}	WorkspaceUpdateRes
// @Failure		default						{object}	Error
// @Security	Authorization
func UseWorkspaceUpdate(service *portal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req WorkspaceUpdateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.ErrorString("PORTAl.WORKSPACE.UPDATE.DECODE.ERROR"))
			return
		}

		in := &usecase.WorkspaceUpdateIn{
			Id:   chi.URLParam(r, "id"),
			Name: req.Name,
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Workspace().Update(r.Context(), in)
		if err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		httpxwriter.Ok(w, &WorkspaceUpdateRes{
			Id:        out.Id,
			CreatedAt: out.CreatedAt,
			UpdatedAt: out.UpdatedAt,
			OwnerId:   out.OwnerId,
			Name:      out.Name,
			Tier:      out.Tier,
		})
	}
}

type WorkspaceUpdateReq struct {
	Name string `json:"name" example:"anothr workspace name"`
} // @name WorkspaceUpdateReq

type WorkspaceUpdateRes struct {
	Id        string `json:"id" example:"ws_2nR9p4W6UmUieJMLIf7ilbXBIRR"`
	CreatedAt int64  `json:"created_at" example:"1728925200000"`
	UpdatedAt int64  `json:"updated_at" example:"1728925200000"`
	OwnerId   string `json:"owner_id" example:"admin"`
	Name      string `json:"name" example:"anothr workspace name"`
	Tier      string `json:"tier" example:"default"`
} // @name WorkspaceUpdateRes
