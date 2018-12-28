package context

import (
	"github.com/jansorg/gotime/gotime/query"
	"github.com/jansorg/gotime/gotime/store"
)

type GoTimeContext struct {
	Store      store.Store
	Query      query.StoreQuery
	JsonOutput bool
}
