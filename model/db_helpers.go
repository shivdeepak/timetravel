package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDb() *gorm.DB {
	if db != nil {
		return db
	}

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
