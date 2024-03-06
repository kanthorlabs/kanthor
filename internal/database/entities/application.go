package entities

import (
	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/validator"
)

type Application struct {
	Auditable

	WsId string
	Name string
}

func (entity *Application) SetId() {
	entity.Id = idx.New(IdNsApp)
}

func (entity *Application) TableName() string {
	return TableApp
}

func (entity *Application) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("id", entity.Id, IdNsApp),
		validator.StringStartsWith("ws_id", entity.WsId, IdNsWs),
		validator.StringRequired("name", entity.Name),
	)
}
