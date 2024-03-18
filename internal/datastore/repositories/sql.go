package repositories

import (
	"sync"

	"github.com/kanthorlabs/common/persistence/datastore"
	"gorm.io/gorm"
)

type sqlrepos struct {
	ds datastore.Datastore

	message  *sqlmsg
	request  *sqlreq
	response *sqlres

	mu sync.RWMutex
}

func (repos *sqlrepos) Message() Message {
	repos.mu.Lock()
	defer repos.mu.Unlock()

	if repos.message == nil {
		repos.message = &sqlmsg{client: repos.ds.Client().(*gorm.DB)}
	}

	return repos.message
}

func (repos *sqlrepos) Request() Request {
	repos.mu.Lock()
	defer repos.mu.Unlock()

	if repos.request == nil {
		repos.request = &sqlreq{client: repos.ds.Client().(*gorm.DB)}
	}

	return repos.request
}

func (repos *sqlrepos) Response() Response {
	repos.mu.Lock()
	defer repos.mu.Unlock()

	if repos.response == nil {
		repos.response = &sqlres{client: repos.ds.Client().(*gorm.DB)}
	}

	return repos.response
}
