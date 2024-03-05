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
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		httpxwriter.Ok(w, &WorkspaceCreateRes{
			Id:        out.Id,
			CreatedAt: out.CreatedAt,
			UpdatedAt: out.UpdatedAt,
			OwnerId:   out.OwnerId,
			Name:      out.Name,
			Tier:      out.Tier,
		})
	}
}

type WorkspaceCreateReq struct {
	Name string `json:"name" example:"main workspace"`
} // @name WorkspaceCreateReq

type WorkspaceCreateRes struct {
	Id        string `json:"id" example:"ws_2nR9p4W6UmUieJMLIf7ilbXBIRR"`
	CreatedAt int64  `json:"created_at" example:"1728925200000"`
	UpdatedAt int64  `json:"updated_at" example:"1728925200000"`
	OwnerId   string `json:"owner_id" example:"admin"`
	Name      string `json:"name" example:"main workspace"`
	Tier      string `json:"tier" example:"default"`
} // @name WorkspaceCreateRes
