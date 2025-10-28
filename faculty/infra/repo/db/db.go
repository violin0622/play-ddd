package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDB() (db *gorm.DB, err error) {
	return gorm.Open(
		sqlite.Open(`:memory:`),
		&gorm.Config{})
}
