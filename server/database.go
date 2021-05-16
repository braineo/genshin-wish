package server

import (
	"github.com/braineo/genshin-wish/parser"

	"github.com/jinzhu/gorm"
	// database driver
	_ "github.com/mattn/go-sqlite3"
)

// GachaLog is Gacha log table record

// InitDB check if database is initiailized
func InitDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "genshin.db")
	if err != nil {
		panic(err)
	}

	db.LogMode(true)

	db.AutoMigrate(&WishLog{}, &parser.GachaConfig{}, &parser.GachaItem{})

	return db
}
