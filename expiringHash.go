// expiringHash project expiringHash.go
package expiringHash

import (
	"sync"
	"sync/atomic"
	"time"
)

type (
	ExpiringHash struct {
		sync.RWMutex
		data        map[string]*expiringItem
		puts        int64
		expirations int64
		getHits     int64
		getMisses   int64
	}

	ExpiringHashStats struct {
		Puts      int64
		GetHits   int64
		GetMisses int64
		Expired   int64
	}

	expiringItem struct {
		value interface{}
		timer *time.Timer
	}
)

// New creates a new initialized struct and returns a pointer.
func New() *ExpiringHash {
	return &ExpiringHash{data: make(map[string]*expiringItem)}
}

// Stats returns usage statistics.
func (eh *ExpiringHash) Stats() ExpiringHashStats {
	eh.RLock()
	defer eh.RLock()
	return ExpiringHashStats{Puts: eh.puts, GetHits: eh.getHits, GetMisses: eh.getMisses, Expired: eh.expirations}
}

// Len returns the number of items.
func (eh *ExpiringHash) Len() int {
	eh.RLock()
	defer eh.RUnlock()
	return len(eh.data)
}

// Put inserts the key/value and sets the key's time-to-live.
func (eh *ExpiringHash) Put(key string, value interface{}, ttl time.Duration) {
	eh.Lock()
	defer eh.Unlock()
	if item, found := eh.data[key]; found {
		item.timer.Reset(ttl)
		item.value = value
		return
	}
	eh.puts = atomic.AddInt64(&eh.puts, 1)
	eh.data[key] = &expiringItem{value: value, timer: time.AfterFunc(ttl, func() { eh.expirations = atomic.AddInt64(&eh.expirations, 1); eh.Del(key) })}
}

// Get returns the value for the key and a boolean indicating if it was found.
func (eh *ExpiringHash) Get(key string) (value interface{}, found bool) {
	eh.RLock()
	defer eh.RUnlock()
	item, found := eh.data[key]
	if found {
		eh.getHits = atomic.AddInt64(&eh.getHits, 1)
		return item.value, true
	}
	eh.getMisses = atomic.AddInt64(&eh.getMisses, 1)
	return "", false
}

// Del deletes an entry.
func (eh *ExpiringHash) Del(key string) {
	eh.Lock()
	defer eh.Unlock()
	if eh.data[key].timer != nil {
		eh.data[key].timer.Stop()
	}
	delete(eh.data, key)
}
