package metricsx

import (
	"expvar"
	"sync"
)

var (
	registryMu sync.Mutex
	registry   = expvar.NewMap("starflow_scheduler")
)

func Inc(name string) {
	Add(name, 1)
}

func Add(name string, delta int64) {
	counter := getCounter(name)
	counter.Add(delta)
}

func Set(name string, value int64) {
	counter := getCounter(name)
	counter.Set(value)
}

func getCounter(name string) *expvar.Int {
	registryMu.Lock()
	defer registryMu.Unlock()
	if existing := registry.Get(name); existing != nil {
		if counter, ok := existing.(*expvar.Int); ok {
			return counter
		}
	}
	counter := new(expvar.Int)
	registry.Set(name, counter)
	return counter
}
