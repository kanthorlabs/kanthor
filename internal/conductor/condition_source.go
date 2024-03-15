package conductor

import (
	"github.com/kanthorlabs/common/validator"
	dbentities "github.com/kanthorlabs/kanthor/internal/database/entities"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

// ConditionSource is a struct that represents a condition source
type ConditionSource struct {
	Source string
}

func (cs *ConditionSource) Validate() error {
	return validator.StringOneOf("source", cs.Source, []string{dbentities.RouteSourceType})()
}

func (cs *ConditionSource) ExtractMessage(msg *dsentities.Message) string {
	if cs.Source == dbentities.RouteSourceType {
		return msg.Type
	}

	return ""
}
