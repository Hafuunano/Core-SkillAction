package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Store is a DB-only store: all operations hit SQLite. Use cache package for in-memory layer.
type Store struct {
	db *gorm.DB
}

// OpenDB opens SQLite at dbPath and returns a Store. No tables are created until Migrate is called.
func OpenDB(dbPath string) (*Store, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// Migrate runs AutoMigrate for the given models. Call once at startup with all models this store needs.
func (s *Store) Migrate(models ...interface{}) error {
	if len(models) == 0 {
		return nil
	}
	return s.db.AutoMigrate(models...)
}

// DB returns the underlying GORM DB for raw queries or custom CRUD.
func (s *Store) DB() *gorm.DB {
	return s.db
}

// Get reads one entry by key from DB only. Returns ("", false, nil) if not found.
func (s *Store) Get(key string) (value string, found bool, err error) {
	var row Entry
	result := s.db.Where("key = ?", key).First(&row)
	if result.Error == gorm.ErrRecordNotFound {
		return "", false, nil
	}
	if result.Error != nil {
		return "", false, result.Error
	}
	return row.Value, true, nil
}

// Set creates or updates an entry in DB only.
func (s *Store) Set(key, value string) error {
	row := Entry{Key: key, Value: value}
	return s.db.Save(&row).Error
}

// Delete removes the entry by key in DB only (hard delete so the same key can be re-inserted).
func (s *Store) Delete(key string) error {
	return s.db.Unscoped().Where("key = ?", key).Delete(&Entry{}).Error
}

// List returns all entries from DB.
func (s *Store) List() ([]Entry, error) {
	var rows []Entry
	if err := s.db.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
