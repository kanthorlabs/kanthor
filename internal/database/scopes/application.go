package scopes

import (
	"fmt"

	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"gorm.io/gorm"
)

func UseApp(wsId string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		joinstm := fmt.Sprintf(
			"JOIN %s ON %s.id = %s.ws_id",
			entities.TableWs, entities.TableWs, entities.TableApp,
		)
		wherestm := fmt.Sprintf("%s.id = ?", entities.TableWs)

		return tx.Joins(joinstm).Where(wherestm, wsId)
	}
}
