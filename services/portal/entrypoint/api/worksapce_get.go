package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

// UseWorkspaceGet
// @Tags			workspace
// @Router		/workspace/{id}		[get]
// @Param			id								path			string						true	"workspace id"
// @Success		200								{object}	WorkspaceGetRes
// @Failure		default						{object}	Error
// @Security	Authorization
func UseWorkspaceGet(service *portal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := &usecase.WorkspaceGetIn{
			Id: chi.URLParam(r, "id"),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Workspace().Get(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &WorkspaceGetRes{Workspace: &Workspace{}}
		res.Map(out.Workspace)
		httpxwriter.Ok(w, res)
	}
}

type WorkspaceGetRes struct {
	*Workspace
} // @name WorkspaceGetRes
