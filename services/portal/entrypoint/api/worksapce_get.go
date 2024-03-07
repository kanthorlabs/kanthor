package api

import (
	"net/http"

	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
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
		ws := r.Context().Value(CtxWorksspace).(*entities.Workspace)
		httpxwriter.Ok(w, ws)
	}
}

type WorkspaceGetRes struct {
	Id        string `json:"id" example:"ws_2nR9p4W6UmUieJMLIf7ilbXBIRR"`
	CreatedAt int64  `json:"created_at" example:"1728925200000"`
	UpdatedAt int64  `json:"updated_at" example:"1728925200000"`
	OwnerId   string `json:"owner_id" example:"admin"`
	Name      string `json:"name" example:"main workspace"`
	Tier      string `json:"tier" example:"default"`
} // @name WorkspaceGetRes
