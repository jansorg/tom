package context

import (
	"github.com/jansorg/gotime/gotime/query"
	"github.com/jansorg/gotime/gotime/store"
)

type GoTimeContext struct {
	Store       store.Store
	StoreHelper *store.StoreHelper
	Query       query.StoreQuery
	JsonOutput  bool
}
