package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/braineo/genshin-wish/server"
	"github.com/op/go-logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var log = logging.MustGetLogger("item-sync")

type Item struct {
	Name   string `json:"name"`
	Rarity string `json:rarity`
}

func main() {
	formatter := logging.MustStringFormatter("%{color}%{time:2006-01-02T15:04:05.000000-07:00} %{module} [%{level}] <%{pid}> [%{shortfile} %{shortfunc}] %{message}%{color:reset}")
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), formatter))
	if level, err := logging.LogLevel("Info"); err == nil {
		logging.SetLevel(level, "")
	}

	db, err := gorm.Open(sqlite.Open("../../genshin.db"), &gorm.Config{})
	if err != nil {
		log.Error(err)
	}
	db.AutoMigrate(&server.GachaItem{})

	types := []string{"character", "weapon"}
	for _, itemType := range types {
		filepath.Walk(fmt.Sprintf("../../genshin-db/src/data/English/%ss", itemType), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Error(err.Error())
			}

			if info.IsDir() {
				return nil
			}
			filename := info.Name()
			log.Infof("Adding %s", info.Name())
			id := filepath.Base(filename)[0 : len(filename)-len(filepath.Ext(filename))]

			f, err := os.Open(filepath.Join(path))
			if err != nil {
				log.Error(err.Error())
			}
			byteValue, err := ioutil.ReadAll(f)
			if err != nil {
				log.Error(err.Error())
			}
			var item Item
			if err = json.Unmarshal(byteValue, &item); err != nil {
				log.Error(err.Error())
			}
			db.FirstOrCreate(&server.GachaItem{
				ID:     id,
				Type:   itemType,
				Name:   item.Name,
				Rarity: item.Rarity,
			})
			return nil
		})
	}

}
