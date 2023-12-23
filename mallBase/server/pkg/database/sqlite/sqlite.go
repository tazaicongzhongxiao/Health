package sqlite

import (
	"MyTestMall/mallBase/basics/pkg/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	Database struct {
		ID        string `gorm:"primary_key;column:id;" json:"id" form:"id"`
		CreatedAt int    `gorm:"column:created_at;index:created_at" json:"created_at"`
		UpdatedAt int    `gorm:"column:updated_at;index:updated_at" json:"updated_at"`
	}
)

var (
	sqLite *gorm.DB
	err    error
)

func Start(dbName string) {
	sqLite, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Warn("database connect error, you can't use sqllite support", err.Error())
	}
}

func Get() *gorm.DB {
	return sqLite
}
