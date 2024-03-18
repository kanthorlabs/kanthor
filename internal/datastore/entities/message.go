package entities

import (
	"time"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/validator"
)

type Message struct {
	Timeseries

	Tier     string
	AppId    string
	Type     string
	Body     string
	Metadata *safe.Metadata
}

func (entity *Message) TableName() string {
	return TableMsg
}

func (entity *Message) SetId() {
	entity.Id = idx.New(IdNsMsg)
}

func (entity *Message) SetTimeseries(now time.Time) {
	if entity.CreatedAt == 0 {
		entity.CreatedAt = now.UnixMilli()
	}
}

func (entity *Message) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("id", entity.Id, IdNsMsg),
		validator.StringRequired("tier", entity.Tier),
		validator.StringRequired("app_id", entity.AppId),
		validator.StringAlphaNumericUnderscoreDot("type", entity.Type),
		validator.StringRequired("body", entity.Body),
	)
}
