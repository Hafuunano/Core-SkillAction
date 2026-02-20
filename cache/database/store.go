package database

import (
	coredb "github.com/Hafuunano/Core-SkillAction/database"
	"gorm.io/gorm"
)

// Store is the database-backed cache: read from memory, write to DB then memory (fastSync). Uses database.
type Store struct {
	db  *gorm.DB
	mem *mem
}

// New builds a database-backed cache Store. Call LoadInMemory() once at startup, then Get/Set/Delete/List.
func New(db *gorm.DB) *Store {
	return &Store{db: db, mem: newMem()}
}

// NewDBCache is an alias for New. Use for database-backed cache.
func NewDBCache(db *gorm.DB) *Store {
	return New(db)
}

// LoadInMemory loads all entries from DB into memory. Call once at startup.
func (s *Store) LoadInMemory() error {
	return s.mem.loadInMemory(s.db)
}

// Get returns value by key from memory first; on miss, loads from DB and backfills memory.
func (s *Store) Get(key string) (value string, found bool, err error) {
	if e, ok := s.mem.getByKey(key); ok {
		return e.Value, true, nil
	}
	var row coredb.Entry
	result := s.db.Where("key = ?", key).First(&row)
	if result.Error == gorm.ErrRecordNotFound {
		return "", false, nil
	}
	if result.Error != nil {
		return "", false, result.Error
	}
	s.mem.setInMemoryEntry(row)
	return row.Value, true, nil
}

// Set writes to DB then updates memory (fastSync write path).
func (s *Store) Set(key, value string) error {
	row := coredb.Entry{Key: key, Value: value}
	if err := s.db.Save(&row).Error; err != nil {
		return err
	}
	s.mem.setInMemoryEntry(row)
	return nil
}

// Delete removes from DB then from memory.
func (s *Store) Delete(key string) error {
	result := s.db.Where("key = ?", key).Delete(&coredb.Entry{})
	if result.Error != nil {
		return result.Error
	}
	s.mem.deleteInMemory(key)
	return nil
}

// GetByKey returns the full Entry from memory only (no DB backfill). Use Get for read-through.
func (s *Store) GetByKey(key string) (coredb.Entry, bool) {
	return s.mem.getByKey(key)
}

// List returns all entries from memory.
func (s *Store) List() []coredb.Entry {
	return s.mem.list()
}
