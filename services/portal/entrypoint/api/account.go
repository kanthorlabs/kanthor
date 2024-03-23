package api

import (
	"github.com/go-chi/chi/v5"
	ppentities "github.com/kanthorlabs/common/passport/entities"
)

func RegisterAccountRoutes(router chi.Router, service *portal) {
	router.Route("/account", func(sr chi.Router) {
		sr.Get("/", UseAccountGet(service))
	})
}

type Account struct {
	*ppentities.Account
} // @name Account

func (acc *Account) Map(entity *ppentities.Account) {
	acc.Account = entity.Censor()
}
