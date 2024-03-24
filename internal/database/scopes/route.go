package scopes

import (
	"fmt"

	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"gorm.io/gorm"
)

func UseRt(wsId string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		joinEp := fmt.Sprintf("JOIN %s ON %s.id = %s.ep_id", entities.TableEp, entities.TableEp, entities.TableRt)
		joinapp := fmt.Sprintf("JOIN %s ON %s.id = %s.app_id", entities.TableApp, entities.TableApp, entities.TableEp)
		joinws := fmt.Sprintf("JOIN %s ON %s.id = %s.ws_id", entities.TableWs, entities.TableWs, entities.TableApp)
		wherestm := fmt.Sprintf("%s.id = ?", entities.TableWs)

		return tx.Joins(joinEp).Joins(joinapp).Joins(joinws).Where(wherestm, wsId)
	}
}
