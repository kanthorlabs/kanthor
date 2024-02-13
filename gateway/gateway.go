package gateway

import (
	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/common/timer"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
)

type Query struct {
	// common
	Search string `json:"_q" form:"_q"`
	Limit  int    `json:"_limit" form:"_limit"`

	// paging
	Page int      `json:"_page" form:"_page"`
	Id   []string `json:"id" form:"id"`

	// scanning
	Start int64 `json:"_start" form:"_start"`
	End   int64 `json:"_end" form:"_end"`
}

func (query *Query) PagingQuery() (*database.PagingQuery, error) {
	err := validator.Validate(
		validator.StringLen("GATEWAY.QUERY.SEARCH", query.Search, 0, database.SearchMaxChar),
		validator.NumberLessThanOrEqual("GATEWAY.QUERY.LIMIT", query.Limit, database.LimitMax),
		validator.NumberLessThanOrEqual("GATEWAY.QUERY.PAGE", query.Page, database.PageMax),
		validator.SliceMaxLength("GATEWAY.QUERY.ID", query.Id, database.PageMax),
	)
	if err != nil {
		return nil, err
	}

	paging := &database.PagingQuery{
		Search: query.Search,
		Limit:  utils.Min(utils.Max(query.Limit, database.LimitMin), database.LimitMax),
		Page:   utils.Min(utils.Max(query.Page, database.PageMin), database.PageMax),
		Ids:    query.Id,
	}

	return paging, nil
}

func (query *Query) ScanningQuery(timer timer.Timer) (*datastore.ScanningQuery, error) {
	err := validator.Validate(
		validator.StringLen("GATEWAY.QUERY.SEARCH.", query.Search, 0, datastore.SearchMaxChar),
		validator.NumberLessThanOrEqual("GATEWAY.QUERY.LIMIT.", query.Limit, datastore.SizeMax),
		validator.NumberLessThanOrEqual("GATEWAY.QUERY.START", query.Start, query.End),
		validator.NumberLessThanOrEqual("GATEWAY.QUERY.START", query.End, timer.Now().UnixMilli()),
	)
	if err != nil {
		return nil, err
	}

	scanning := &datastore.ScanningQuery{
		Search: query.Search,
		Size:   utils.Min(utils.Max(query.Limit, datastore.SizeMin), datastore.SizeMax),
		From:   timer.UnixMilli(query.Start),
		To:     timer.UnixMilli(query.End),
	}

	return scanning, nil
}
