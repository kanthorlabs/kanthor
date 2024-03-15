package entities

import (
	"time"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/safe"
)

type Request struct {
	Timeseries

	// message properties
	MsgId    string
	Tier     string
	AppId    string
	Type     string
	Metadata *safe.Metadata
	Body     string

	// endpoint properties
	EpId string
	// HTTP: POST/PUT/PATCH
	Method  string
	Uri     string
	Headers *safe.Metadata
}

func (entity *Request) SetId() {
	entity.Id = idx.New(IdNsReq)
}

func (entity *Request) SetTimeseries(now time.Time) {
	if entity.CreatedAt == 0 {
		entity.CreatedAt = now.UnixMilli()
	}
}
