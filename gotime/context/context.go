package context

import "../store"

type GoTimeContext struct {
	Store      store.Store
	JsonOutput bool
}
