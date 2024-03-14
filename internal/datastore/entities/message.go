package entities

import "github.com/kanthorlabs/common/safe"

type Message struct {
	Timeseries

	Tier     string
	AppId    string
	Type     string
	Metadata *safe.Metadata
	Headers  *safe.Metadata
	Body     string
}
