package sqly

import (
	"github.com/jmoiron/sqlx"
	"sync"
)

var scriptsCache = sync.Map{}
var databasesCache = sync.Map{}

type databasesCacheValue struct {
	db   *sqlx.DB
	name string
	ping bool
}

type scriptsCacheValue struct {
	name     string
	content  string
	database string
	path     string
}

func lintDatabasesCache(cachedDatabase []string) {
	f := func(key interface{}, value interface{}) bool {
		keyStr := key.(string)
		if !contains(cachedDatabase, keyStr) {
			databasesCache.Delete(key)
		}
		return true
	}
	databasesCache.Range(f)
}
