// Package core provides a unified Services facade (DB, Cache, Timer) for Lucy to inject into protocol.Context.
package core

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	cachedb "github.com/Hafuunano/Core-SkillAction/cache/database"
	"github.com/Hafuunano/Core-SkillAction/database"
	"github.com/Hafuunano/Core-SkillAction/timer"
)

const defaultDBPath = "data/skill_action.db"

var (
	defaultOnce  sync.Once
	defaultCache *cachedb.Store
)

// ServicesOptions configures NewServices. DBPath is required for DB and for DB-backed cache when EnableDBCache is true.
type ServicesOptions struct {
	// DBPath is the SQLite file path. Required when using DB or DB-backed cache.
	DBPath string
	// EnableDBCache when true creates a database-backed cache and calls LoadInMemory at startup.
	EnableDBCache bool
	// TimerTTL is the default TTL for the Timer store. If zero, 10 minutes is used.
	TimerTTL time.Duration
}

const defaultTimerTTL = 10 * time.Minute

// Services holds injectable DB, Cache, and Timer. Lucy calls NewServices once at startup and passes *Services when constructing each message Context.
type Services struct {
	// DB is the SQLite store. Never nil when NewServices succeeds with non-empty DBPath.
	DB *database.Store
	// Cache is the database-backed cache when EnableDBCache was true; nil otherwise. Callers can check Cache != nil before use.
	Cache *cachedb.Store
	// Timer is a TTL key-value store (string -> any). Never nil.
	Timer *timer.Store[string, any]
}

// NewServices creates DB (and optionally DB-backed cache) and a default Timer store. Migrate and cache LoadInMemory are called inside.
// Pass the returned *Services to Lucy so it can inject into protocol.Context when handling each event.
func NewServices(opts ServicesOptions) (*Services, error) {
	ttl := opts.TimerTTL
	if ttl <= 0 {
		ttl = defaultTimerTTL
	}

	s := &Services{
		Timer: timer.NewStore[string, any](ttl),
	}

	if opts.DBPath != "" {
		dbStore, err := database.OpenDB(opts.DBPath)
		if err != nil {
			return nil, err
		}
		if err := dbStore.Migrate(&database.Entry{}); err != nil {
			return nil, err
		}
		s.DB = dbStore

		if opts.EnableDBCache {
			cacheStore := cachedb.New(dbStore.DB())
			if err := cacheStore.LoadInMemory(); err != nil {
				return nil, err
			}
			s.Cache = cacheStore
		}
	}

	return s, nil
}

// DefaultCache returns a lazily-created database-backed cache Store using default path (data/skill_action.db).
// Plugins that need a store can use this when the host has not called SetStore; the first call creates the DB and cache.
func DefaultCache() *cachedb.Store {
	defaultOnce.Do(func() {
		_ = os.MkdirAll(filepath.Dir(defaultDBPath), 0755)
		svc, err := NewServices(ServicesOptions{
			DBPath:         defaultDBPath,
			EnableDBCache:  true,
		})
		if err != nil {
			return
		}
		if svc != nil {
			defaultCache = svc.Cache
		}
	})
	return defaultCache
}
