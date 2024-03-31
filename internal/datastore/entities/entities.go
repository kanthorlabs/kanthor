package entities

import (
	"github.com/kanthorlabs/common/project"
)

var (
	IdNsMsg = "msg"
	IdNsReq = "req"
	IdNsRes = "res"

	TableMsg = project.Name("message")
	TableReq = project.Name("request")
	TableRes = project.Name("response")
)

// Timeseries is a base entity that contains timeseries data such as id and created_at
type Timeseries struct {
	Id string
	// I didn't find a way to disable automatic fields modify yet
	// so, I use a tag to disable this feature here
	// but, we should keep our entities stateless if we can
	CreatedAt int64 `gorm:"autoCreateTime:false"`
}

type MessagePk struct {
	AppId string
	Id    string
}

type RequestPk struct {
	EpId  string
	MsgId string
}

type ResponsePk struct {
	EpId  string
	MsgId string
}
