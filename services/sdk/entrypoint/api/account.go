package api

import (
	"github.com/go-chi/chi/v5"
	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/safe"
)

func RegisterAccountRoutes(router chi.Router, service *sdk) {
	router.Route("/account", func(sr chi.Router) {
		sr.Get("/", UseAccountGet(service))
	})
}

type Account struct {
	Username string         `json:"username"`
	Name     string         `json:"name"`
	Metadata *safe.Metadata `json:"metadata"`

	CreatedAt     int64 `json:"created_at"`
	UpdatedAt     int64 `json:"updated_at" `
	DeactivatedAt int64 `json:"deactivated_at"`
} // @name Account

func (acc *Account) Map(entity *ppentities.Account) {
	acc.Username = entity.Username
	acc.Name = entity.Name
	acc.Metadata = &safe.Metadata{}
	acc.Metadata.Merge(entity.Metadata)
	acc.CreatedAt = entity.CreatedAt
	acc.UpdatedAt = entity.UpdatedAt
	acc.DeactivatedAt = entity.DeactivatedAt
}
