package entities

import "github.com/kanthorlabs/common/project"

var (
	TableWs  = project.Name("workspace")
	TableWsc = project.Name("workspace_credentials")
	TableApp = project.Name("application")
	TableEp  = project.Name("endpoint")
	TableEpr = project.Name("endpoint_rule")
	TableMsg = project.Name("message")
	TableReq = project.Name("request")
	TableRes = project.Name("response")
	TableAtt = project.Name("attempt")
)
