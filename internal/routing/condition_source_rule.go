package routing

import "github.com/kanthorlabs/kanthor/internal/entities"

func ConditionSource(rule *entities.EndpointRule, msg *entities.Message) string {
	if rule.ConditionSource == ConditionSourceType {
		return msg.Type
	}
	if rule.ConditionSource == ConditionSourceAppId {
		return msg.AppId
	}
	return ""
}
