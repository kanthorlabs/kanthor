package entities

import (
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

type EndpointRule struct {
	Auditable

	EpId string
	Name string

	Priority int32
	// the logic of not-false is true should be used here
	// to guarantee default all rule will be on include mode
	Exclusionary bool

	// examples
	//  - tag
	ConditionSource string
	// examples:
	// 	- any::
	// 	- equal::orders.paid
	// 	- prefix::orders.
	ConditionExpression string
}

func (entity *EndpointRule) SetId() {
	entity.Id = idx.New(IdNsEpr)
}

func (entity *EndpointRule) TableName() string {
	return TableEpr
}

func (entity *EndpointRule) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("id", entity.Id, IdNsEpr),
		validator.StringStartsWith("ep_id", entity.EpId, IdNsEp),
		validator.StringRequired("name", entity.Name),
		validator.NumberGreaterThan("priority", entity.Priority, 0),
		validator.StringOneOf("condition_source", entity.ConditionSource, []string{"tag"}),
		validator.StringRequired("condition_expression", entity.ConditionExpression),
	)
}
