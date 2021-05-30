package server

import (
	"time"

	"gorm.io/gorm"
)

type WishLog struct {
	gorm.Model
	ID        string `json:"id"` // id for pagination
	GachaType string `json:"gachaType"`
	UserID    string `json:"userId"`
	User      User
	// Pull time, should convert from time string to unix timestamp
	Time      time.Time `json:"time"`
	ItemID    string    `json:"itemId"`
	Item      GachaItem
	PityStar4 int `json:"pityStar4"`
	PityStar5 int `json:"pityStar5"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GachaItem struct {
	gorm.Model
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Rarity string `json:"rarity"`
}

type ItemCount struct {
	ItemID string `json:"itemId"`
	Item   GachaItem
	Count  int
}

