// Package database provides cache backed by Gorm and in-memory fastSync. Uses database.
package database

import (
	"sync"

	coredb "github.com/Hafuunano/Core-SkillAction/database"
	"gorm.io/gorm"
)

// mem holds in-memory key-value copy of coredb.Entry for fast read.
type mem struct {
	mu   sync.RWMutex
	data map[string]coredb.Entry
}

func newMem() *mem {
	return &mem{data: make(map[string]coredb.Entry)}
}

func (m *mem) loadInMemory(db *gorm.DB) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var rows []coredb.Entry
	if err := db.Find(&rows).Error; err != nil {
		return err
	}
	m.data = make(map[string]coredb.Entry, len(rows))
	for _, e := range rows {
		m.data[e.Key] = e
	}
	return nil
}

func (m *mem) getByKey(key string) (coredb.Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.data[key]
	return e, ok
}

func (m *mem) setInMemoryEntry(e coredb.Entry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[e.Key] = e
}

func (m *mem) deleteInMemory(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *mem) list() []coredb.Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]coredb.Entry, 0, len(m.data))
	for _, e := range m.data {
		out = append(out, e)
	}
	return out
}
