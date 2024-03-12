package api

import (
	"net/http"

	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

// UseWorkspaceList
// @Tags			workspace
// @Router		/workspace				[get]
// @Success		200								{object}	WorkspaceListRes
// @Failure		default						{object}	Error
// @Security	Authorization
func UseWorkspaceList(service *portal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account := r.Context().Value(passport.CtxAccount).(*ppentities.Account)
		in := &usecase.WorkspaceListIn{
			OwnerId: account.Username,
		}
		if err := in.Validate(); err != nil {
			httpxwriter.ErrBadRequest(w, httpxwriter.Error(err))
			return
		}

		out, err := service.uc.Workspace().List(r.Context(), in)
		if err != nil {
			httpxwriter.ErrUnknown(w, httpxwriter.Error(err))
			return
		}

		res := &WorkspaceListRes{Data: make([]*Workspace, len(out.Data))}
		for i := range out.Data {
			res.Data[i] = &Workspace{}
			res.Data[i].Map(out.Data[i])
		}
		httpxwriter.Ok(w, res)
	}
}

type WorkspaceListRes struct {
	Data []*Workspace `json:"data"`
} // @name WorkspaceListRes
