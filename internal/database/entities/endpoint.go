package entities

import (
	"fmt"
	"net/http"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/validator"
)

type Endpoint struct {
	Auditable

	SecretKey string

	AppId  string
	Name   string
	Method string
	Uri    string
}

func (entity *Endpoint) SetId() {
	entity.Id = idx.New(IdNsEp)
}

func (entity *Endpoint) TableName() string {
	return TableEp
}

func (entity *Endpoint) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("id", entity.Id, IdNsEp),
		validator.StringLen("secret_key", entity.SecretKey, 16, 32),
		validator.StringStartsWith("app_id", entity.AppId, IdNsApp),
		validator.StringRequired("name", entity.Name),
		validator.StringOneOf("method", entity.Method, []string{http.MethodPost, http.MethodPut}),
		validator.StringUri("uri", entity.Uri),
	)
}

func (entity *Endpoint) PrimaryProp() string {
	return fmt.Sprintf("%s.id", TableEp)
}

func (entity *Endpoint) SearchProps() []string {
	return []string{
		fmt.Sprintf("%s.name", TableEp),
		fmt.Sprintf("%s.method", TableEp),
		fmt.Sprintf("%s.uri", TableEp),
	}
}
