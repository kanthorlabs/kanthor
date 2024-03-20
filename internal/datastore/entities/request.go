package entities

import (
	"net/http"
	"time"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/validator"
)

type Request struct {
	Timeseries

	// message properties
	MsgId string
	Tier  string
	AppId string
	Type  string
	Body  string

	// endpoint properties
	EpId string
	// HTTP: POST/PUT
	Method string
	Uri    string

	Headers  *safe.Metadata
	Metadata *safe.Metadata
}

func (entity *Request) TableName() string {
	return TableReq
}

func (entity *Request) SetId() {
	entity.Id = idx.New(IdNsReq)
}

func (entity *Request) SetTimeseries(now time.Time) {
	if entity.CreatedAt == 0 {
		entity.CreatedAt = now.UnixMilli()
	}
}

func (entity *Request) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("id", entity.Id, IdNsReq),
		validator.StringStartsWith("msg_id", entity.MsgId, IdNsMsg),
		validator.StringRequired("tier", entity.Tier),
		validator.StringRequired("app_id", entity.AppId),
		validator.StringAlphaNumericUnderscoreDot("type", entity.Type),
		validator.StringRequired("body", entity.Body),
		validator.StringRequired("ep_id", entity.EpId),
		validator.StringOneOf("method", entity.Method, []string{http.MethodPost, http.MethodPut}),
		validator.StringUri("uri", entity.Uri),
	)
}
