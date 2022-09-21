package mysql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

var (
	dbs  = map[string]*sql.DB{}
	dbRW sync.RWMutex
)

func getDB(name string, opts ...Option) (*sql.DB, error) {
	dbRW.RLock()
	db, ok := dbs[name]
	dbRW.RUnlock()
	if ok {
		return db, nil
	}

	cfg := getServiceConfig(name)
	for _, opt := range opts {
		opt(&cfg)
	}

	return newDB(&cfg)
}

func newDB(cfg *serviceConfig) (*sql.DB, error) {
	dbRW.Lock()
	defer dbRW.Unlock()

	db, ok := dbs[cfg.Name]
	if ok {
		return db, nil
	}

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("mysql open fail. error:%v", err)
	}

	if cfg.MaxIdle > 0 {
		db.SetMaxIdleConns(cfg.MaxIdle)
	}

	if cfg.MaxOpen > 0 {
		db.SetMaxOpenConns(cfg.MaxOpen)
	}

	if cfg.MaxIdleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Millisecond)
	}

	dbs[cfg.Name] = db
	return db, nil
}
