package entities

import (
	"time"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/safe"
)

type Message struct {
	Timeseries

	Tier     string
	AppId    string
	Type     string
	Body     string
	Metadata *safe.Metadata
}

func (entity *Message) SetId() {
	entity.Id = idx.New(IdNsMsg)
}

func (entity *Message) TableName() string {
	return TableMsg
}

func (entity *Message) SetTimeseries(now time.Time) {
	if entity.CreatedAt == 0 {
		entity.CreatedAt = now.UnixMilli()
	}
}
