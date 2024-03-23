package api

import (
	"net/http"

	httpxwriter "github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	ppentities "github.com/kanthorlabs/common/passport/entities"
)

// UseAccountGet
// @Tags			account
// @Router		/account				[get]
// @Success		200							{object}	AccountGetRes
// @Failure		default					{object}	Error
// @Security	Authorization
func UseAccountGet(service *portal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account := r.Context().Value(passport.CtxAccount).(*ppentities.Account)

		res := &AccountGetRes{Account: &Account{}}
		res.Map(account)
		httpxwriter.Ok(w, res)
	}
}

type AccountGetRes struct {
	*Account
} // @name AccountGetRes
