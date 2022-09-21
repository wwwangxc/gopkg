package orm

import (
	"gorm.io/gorm"
)

// NewGORM new GORM DB
func NewGORM(name string, opts ...GORMOption) (*gorm.DB, error) {
	return getGORMDB(name, opts...)
}
