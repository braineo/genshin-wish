package server

import (
	"github.com/braineo/genshin-wish/parser"

	"github.com/jinzhu/gorm"
	// database driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("cli")

// GachaLog is Gacha log table record

// InitDB check if database is initiailized
func InitDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "genshin.db")
	if err != nil {
		panic(err)
	}

	db.LogMode(true)

	db.AutoMigrate(&parser.GachaLog{}, &parser.GachaConfig{})

	return db
}
