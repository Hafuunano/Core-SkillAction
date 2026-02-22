// Package timer
package timer

import (
	"time"

	"github.com/FloatTech/ttl"
)

// Store is a TTL key-value store with create, read, update, delete operations.
// K must be comparable (e.g. string, int64), V can be any type.
type Store[K comparable, V any] struct {
	cache *ttl.Cache[K, V]
}

// NewStore creates a TTL store with the given TTL duration.
// Items expire and are removed automatically after ttl from last Set or Get.
func NewStore[K comparable, V any](ttlDuration time.Duration) *Store[K, V] {
	return &Store[K, V]{
		cache: ttl.NewCache[K, V](ttlDuration),
	}
}

// Set creates or updates a key-value pair (增/改).
// If key exists, value is overwritten and TTL is reset.
func (s *Store[K, V]) Set(key K, val V) {
	s.cache.Set(key, val)
}

// Get reads the value for key (查).
// Returns the zero value of V if key is missing or expired.
func (s *Store[K, V]) Get(key K) (v V) {
	return s.cache.Get(key)
}

// GetOrSet gets the value for key, or sets and returns val if key is missing/expired (查+增).
// Second return is true if the value was already present.
func (s *Store[K, V]) GetOrSet(key K, val V) (v V, existed bool) {
	return s.cache.GetOrSet(key, val)
}

// Delete removes the key (删).
func (s *Store[K, V]) Delete(key K) {
	s.cache.Delete(key)
}

// GetAndDelete gets the value for key and then removes it (查+删).
// Second return is false if key was missing or expired.
func (s *Store[K, V]) GetAndDelete(key K) (v V, deleted bool) {
	return s.cache.GetAndDelete(key)
}

// Touch extends the TTL for key by the given duration.
func (s *Store[K, V]) Touch(key K, extra time.Duration) {
	s.cache.Touch(key, extra)
}

// Range calls f for each non-expired key-value pair. Stops on first error.
func (s *Store[K, V]) Range(f func(K, V) error) error {
	return s.cache.Range(f)
}

// Destroy stops the background GC and clears all items. Do not use the store after Destroy.
func (s *Store[K, V]) Destroy() {
	s.cache.Destroy()
}
