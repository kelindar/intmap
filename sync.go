// Copyright (c) 2021-2025, Roman Atachiants

package intmap

import "sync"

// Sync is a thread-safe, map-like data-structure for int64s
type Sync struct {
	lock sync.RWMutex
	data *Map
}

// NewSync returns a thread-safe map initialized with n spaces and uses the stated fillFactor.
// The map will grow as needed.
func NewSync(size int, fillFactor float64) *Sync {
	return &Sync{
		data: New(size, fillFactor),
	}
}

// Load returns the value stored in the map for a key, or nil if no value is
// present. The ok result indicates whether value was found in the map.
func (m *Sync) Load(key uint32) (value uint32, ok bool) {
	m.lock.RLock()
	value, ok = m.data.Load(key)
	m.lock.RUnlock()
	return
}

// Store sets the value for a key.
func (m *Sync) Store(key, val uint32) {
	m.lock.Lock()
	m.data.Store(key, val)
	m.lock.Unlock()
}

// Delete deletes the value for a key.
func (m *Sync) Delete(key uint32) {
	m.lock.Lock()
	m.data.Delete(key)
	m.lock.Unlock()
}

// Count returns number of key/value pairs in the map.
func (m *Sync) Count() (count int) {
	m.lock.RLock()
	count = m.data.Count()
	m.lock.RUnlock()
	return
}

// LoadOrStore returns the existing value for the key if present. Otherwise, it stores
// and returns the given value returned by the handler. The loaded result is true if the
// value was loaded, false if stored.
func (m *Sync) LoadOrStore(key uint32, fn func() uint32) (value uint32, loaded bool) {
	if value, loaded = m.Load(key); loaded {
		return // fast-path
	}

	// Load or store again, with exclusive lock now
	m.lock.Lock()
	defer m.lock.Unlock()
	if value, loaded = m.data.Load(key); !loaded {
		value = fn()
		m.data.Store(key, value)
	}
	return
}

// Range calls f sequentially for each key and value present in the map. If f
// returns false, range stops the iteration.
func (m *Sync) Range(f func(key, value uint32) bool) {
	m.lock.RLock()
	m.data.Range(f)
	m.lock.RUnlock()
}
