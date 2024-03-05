package entities

type SnaptshotRow struct {
	WsId                   string `json:"ws_id"`
	WsName                 string `json:"ws_name"`
	AppId                  string `json:"app_id"`
	AppName                string `json:"app_name"`
	EpId                   string `json:"ep_id"`
	EpName                 string `json:"ep_name"`
	EpMethod               string `json:"ep_method"`
	EpUri                  string `json:"ep_uri"`
	EprId                  string `json:"epr_id"`
	EprName                string `json:"epr_name"`
	EprPriority            int32  `json:"epr_priority"`
	EprExclusionary        bool   `json:"epr_exclusionary"`
	EprConditionSource     string `json:"epr_condition_source"`
	EprConditionExpression string `json:"epr_condition_expression"`
}
