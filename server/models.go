package server

import (
	"time"

	"github.com/braineo/genshin-wish/parser"
)

type WishLog struct {
	ID        string `json:"id"` // id for pagination
	GachaType string `json:"gacha_type"`
	UserID    string `json:"uid"`
	// Pull time, should convert from time string to unix timestamp
	Time      time.Time `json:"time"`
	ItemID    string
	Item      parser.GachaItem
	PityStar4 int `json:"pityStar4"`
	PityStar5 int `json:"pityStar5"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
