package entities

import (
	"fmt"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/validator"
)

type Workspace struct {
	Auditable

	OwnerId string
	Name    string
	Tier    string
}

func (entity *Workspace) TableName() string {
	return TableWs
}

func (entity *Workspace) SetId() {
	entity.Id = idx.New(IdNsWs)
}

func (entity *Workspace) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("id", entity.Id, IdNsWs),
		validator.StringRequired("owner_id", entity.OwnerId),
		validator.StringRequired("name", entity.Name),
		validator.StringRequired("tier", entity.Tier),
	)
}

func (entity *Workspace) PrimaryProp() string {
	return fmt.Sprintf("%s.id", TableWs)
}

func (entity *Workspace) SearchProps() []string {
	return []string{
		fmt.Sprintf("%s.name", TableWs),
	}
}
