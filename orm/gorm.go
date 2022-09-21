package orm

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/wwwangxc/gopkg/orm/driver"
)

var (
	gormDBMap   = map[string]*gorm.DB{}
	gormDBMapRW sync.RWMutex
)

func getGORMDB(name string, opts ...GORMOption) (*gorm.DB, error) {
	gormDBMapRW.RLock()
	db, ok := gormDBMap[name]
	gormDBMapRW.RUnlock()
	if ok {
		return db, nil
	}

	cfg := getServiceConfig(name)
	for _, opt := range opts {
		opt(&cfg)
	}

	return newGORMDB(&cfg)
}

func newGORMDB(cfg *serviceConfig) (*gorm.DB, error) {
	gormDBMapRW.Lock()
	defer gormDBMapRW.Unlock()

	db, ok := gormDBMap[cfg.Name]
	if ok {
		return db, nil
	}

	d, exist := driver.Get(cfg.Driver)
	if !exist {
		return nil, fmt.Errorf("invalid driver:%s", cfg.Driver)
	}

	var opts []gorm.Option
	if cfg.gormConfig != nil {
		opts = append(opts, cfg.gormConfig)
	}

	db, err := gorm.Open(d.Open(cfg.DSN), opts...)
	if err != nil {
		return nil, fmt.Errorf("gorm open fail. error:%v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB fail. error:%v", err)
	}

	if cfg.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	}

	if cfg.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	}

	if cfg.MaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Millisecond)
	}

	gormDBMap[cfg.Name] = db
	return db, nil
}
