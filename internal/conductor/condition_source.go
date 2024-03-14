package conductor

import (
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

// ConditionSource is a struct that represents a condition source
type ConditionSource struct {
	Source string
}

func (cs *ConditionSource) Validate() error {
	return validator.StringOneOf("source", cs.Source, []string{entities.RouteSourceTag})()
}

func (cs *ConditionSource) ExtractMessage() string {
	return ""
}
