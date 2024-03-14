package entities

import (
	"fmt"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/validator"
)

var RouteSourceTag = "tag"

type Route struct {
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

func (entity *Route) SetId() {
	entity.Id = idx.New(IdNsRt)
}

func (entity *Route) TableName() string {
	return TableRt
}

func (entity *Route) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("id", entity.Id, IdNsRt),
		validator.StringStartsWith("ep_id", entity.EpId, IdNsEp),
		validator.StringRequired("name", entity.Name),
		validator.NumberInRange("priority", entity.Priority, 1, 128),
		validator.StringOneOf("condition_source", entity.ConditionSource, []string{RouteSourceTag}),
		validator.StringRequired("condition_expression", entity.ConditionExpression),
	)
}

func (entity *Route) PrimaryProp() string {
	return fmt.Sprintf("%s.id", TableRt)
}

func (entity *Route) SearchProps() []string {
	return []string{
		fmt.Sprintf("%s.name", TableRt),
	}
}
