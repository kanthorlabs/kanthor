package entities

import (
	"time"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/validator"
)

type Response struct {
	Timeseries

	// message properties
	MsgId string
	Tier  string
	AppId string
	Type  string

	// request properties
	ReqId    string
	EpId     string
	Headers  *safe.Metadata
	Metadata *safe.Metadata

	Status int
	Uri    string
	Body   string
}

func (entity *Response) TableName() string {
	return TableRes
}

func (entity *Response) SetId() {
	entity.Id = idx.New(IdNsRes)
}

func (entity *Response) SetTimeseries(now time.Time) {
	if entity.CreatedAt == 0 {
		entity.CreatedAt = now.UnixMilli()
	}
}

func (entity *Response) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("id", entity.Id, IdNsRes),
		validator.StringStartsWith("msg_id", entity.MsgId, IdNsMsg),
		validator.StringRequired("tier", entity.Tier),
		validator.StringRequired("app_id", entity.AppId),
		validator.StringAlphaNumericUnderscoreHyphenDot("type", entity.Type),
		validator.StringStartsWith("req_id", entity.ReqId, IdNsReq),
		validator.StringRequired("ep_id", entity.EpId),
		validator.StringUri("uri", entity.Uri),
		validator.StringRequired("body", entity.Body),
	)
}
