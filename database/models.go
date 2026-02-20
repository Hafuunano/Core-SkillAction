// Package database provides Gorm+SQLite connection and direct DB CRUD only (no in-memory cache).
package database

import (
	"time"

	"gorm.io/gorm"
)

// Entry is the key-value row stored in SQLite.
type Entry struct {
	ID        uint           `gorm:"primaryKey"`
	Key       string         `gorm:"uniqueIndex;size:512;not null"`
	Value     string         `gorm:"type:text"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName overrides the table name.
func (Entry) TableName() string {
	return "skill_action_entries"
}
