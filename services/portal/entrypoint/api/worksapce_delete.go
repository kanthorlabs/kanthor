package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

// UseWorkspaceDelete
// @Tags			workspace
// @Router		/workspace/{id}		[delete]
// @Param			id								path			string						true	"workspace id"
// @Success		200								{object}	WorkspaceDeleteRes
// @Failure		default						{object}	Error
// @Security	Authorization
func UseWorkspaceDelete(service *portal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := &usecase.WorkspaceDeleteIn{
			Id: chi.URLParam(r, "id"),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Workspace().Delete(r.Context(), in)
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

type WorkspaceDeleteRes struct {
	Id        string `json:"id" example:"ws_2nR9p4W6UmUieJMLIf7ilbXBIRR"`
	CreatedAt int64  `json:"created_at" example:"1728925200000"`
	DeletedAt int64  `json:"updated_at" example:"1728925200000"`
	OwnerId   string `json:"owner_id" example:"admin"`
	Name      string `json:"name" example:"main workspace"`
	Tier      string `json:"tier" example:"default"`
} // @name WorkspaceDeleteRes
