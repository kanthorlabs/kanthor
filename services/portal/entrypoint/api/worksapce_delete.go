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
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &WorkspaceDeleteRes{Workspace: &Workspace{}}
		res.Map(out.Workspace)
		httpxwriter.Ok(w, res)
	}
}

type WorkspaceDeleteRes struct {
	*Workspace
} // @name WorkspaceDeleteRes
