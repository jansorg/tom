package context

import "github.com/jansorg/gotime/gotime/store"

type GoTimeContext struct {
	Store      store.Store
	JsonOutput bool
}
