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

		res := &WorkspaceUpdateRes{Workspace: &Workspace{}}
		res.Map(out.Workspace)
		httpxwriter.Ok(w, res)
	}
}

type WorkspaceUpdateReq struct {
	Name string `json:"name" example:"anthor workspace name"`
} // @name WorkspaceUpdateReq

type WorkspaceUpdateRes struct {
	*Workspace
} // @name WorkspaceUpdateRes
