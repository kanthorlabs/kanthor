package entities

import (
	"time"

	"github.com/kanthorlabs/common/project"
)

var (
	IdNsWs  = "ws"
	IdNsApp = "app"
	IdNsEp  = "ep"
	IdNsRt  = "rt"

	TableWs  = project.Name("workspace")
	TableApp = project.Name("application")
	TableEp  = project.Name("endpoint")
	TableRt  = project.Name("route")
)

// Auditable is a base entity that contains audit fields, such as created_at, updated_at, and modifier
// It can tell us who and when the entity was created and updated
type Auditable struct {
	Id string
	// I didn't find a way to disable automatic fields modify yet
	// so, I use a tag to disable this feature here
	// but, we should keep our entities stateless if we can
	CreatedAt int64 `gorm:"autoCreateTime:false"`
	UpdatedAt int64 `gorm:"autoUpdateTime:false"`
	Modifier  string
}

func (entity *Auditable) SetAuditFacttor(now time.Time, modifier string) {
	if entity.CreatedAt == 0 {
		entity.CreatedAt = now.UnixMilli()
	}
	entity.UpdatedAt = now.UnixMilli()
	entity.Modifier = modifier
}
