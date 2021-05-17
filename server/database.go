package server

import (
	"github.com/braineo/genshin-wish/parser"

	"gorm.io/gorm"
	// database driver
	"gorm.io/driver/sqlite"
)

// GachaLog is Gacha log table record

// InitDB check if database is initiailized
func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("genshin.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&WishLog{}, &parser.GachaConfig{}, &GachaItem{})

	return db
}
