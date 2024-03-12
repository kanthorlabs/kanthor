package api

import (
	"encoding/json"
	"net/http"

	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

// UseWorkspaceCreate
// @Tags			workspace
// @Router		/workspace		[post]
// @Param			request				body			WorkspaceCreateReq	true	"request body"
// @Success		200						{object}	WorkspaceCreateRes
// @Failure		default				{object}	Error
// @Security	Authorization
func UseWorkspaceCreate(service *portal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req WorkspaceCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.ErrorString("PORTAl.WORKSPACE.CREATE.DECODE.ERROR"))
			return
		}

		account := r.Context().Value(passport.CtxAccount).(*ppentities.Account)
		in := &usecase.WorkspaceCreateIn{
			OwnerId: account.Username,
			Name:    req.Name,
			// use default tier
			Tier: project.Tier(),
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Workspace().Create(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &WorkspaceCreateRes{Workspace: &Workspace{}}
		res.Map(out.Workspace)
		httpxwriter.Ok(w, res)
	}
}

type WorkspaceCreateReq struct {
	Name string `json:"name" example:"main workspace"`
} // @name WorkspaceCreateReq

type WorkspaceCreateRes struct {
	*Workspace
} // @name WorkspaceCreateRes
