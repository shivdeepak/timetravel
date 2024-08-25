package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db/dev.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func InitDb() {
	db := GetDb()

	db.AutoMigrate(&Record{})
}
