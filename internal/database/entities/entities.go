package entities

import (
	"time"

	"github.com/kanthorlabs/common/project"
)

var (
	IdNsWs   = "ws"
	IdNsWsc  = "wsc"
	IdNsApp  = "app"
	IdNsEp   = "ep"
	IdNsRt   = "rt"
	IdNsEpsk = "epsk"

	TableWs  = project.Name("workspace")
	TableApp = project.Name("application")
	TableEp  = project.Name("endpoint")
	TableRt  = project.Name("route")
)

type Auditable struct {
	Id string
	// I didn't find a way to disable automatic fields modify yet
	// so, I use a tag to disable this feature here
	// but, we should keep our entities stateless if we can
	CreatedAt int64 `gorm:"autoCreateTime:false"`
	UpdatedAt int64 `gorm:"autoUpdateTime:false"`
}

func (entity *Auditable) SetAuditFacttor(now time.Time) {
	if entity.CreatedAt == 0 {
		entity.CreatedAt = now.UnixMilli()
	}
	entity.UpdatedAt = now.UnixMilli()
}
