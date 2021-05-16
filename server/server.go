package server

import (
	"net/http"
	"time"

	"github.com/braineo/genshin-wish/parser"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/op/go-logging"
)

func handleCors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}

var log = logging.MustGetLogger("server")

type Server struct {
	Engine   *gin.Engine
	Database *gorm.DB
}

func New() Server {

	engine := gin.Default()
	db := InitDB()
	server := Server{
		Engine:   engine,
		Database: db,
	}
	// engine.Use(handleCors())

	v1Route := engine.Group("api/v1")

	{
		v1Route.POST("/log", server.FetchLogs)
		v1Route.GET("/log/:uid", server.GetLogs)

		v1Route.POST("/item", server.FetchGachaItems)
		v1Route.GET("/item", server.GetGachaItems)

		v1Route.POST("/gacha", server.FetchGachaConfigs)
		v1Route.GET("/gacha", server.GetGachaConfigs)

	}
	return server
}

func (server *Server) Run() {

	server.Engine.Run(":8080")
}

func (server *Server) FetchGachaItems(ctx *gin.Context) {
	rawQuery := ctx.PostForm("query")
	p, err := parser.New(rawQuery, parser.WithLanguage(parser.EnUs))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	p.FetchGachaItems()
	for _, item := range p.ItemTable {
		server.Database.FirstOrCreate(item)
	}
}

func (server *Server) GetGachaItems(ctx *gin.Context) {

}

func (server *Server) FetchGachaConfigs(ctx *gin.Context) {

}

func (server *Server) GetGachaConfigs(ctx *gin.Context) {

}

// FetchLogs accept query URL for gacha log to query game server
func (server *Server) FetchLogs(ctx *gin.Context) {
	rawQuery := ctx.PostForm("query")

	p, err := parser.New(rawQuery, parser.WithLanguage(parser.EnUs))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	var configs []parser.GachaConfig
	server.Database.Find(&configs)

	if len(configs) == 0 {
		p.Language = parser.ZhCn
		if err = p.FetchGachaConfigs(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
		p.Language = parser.EnUs
		// TODO insert new gacha configs
		for _, config := range p.Configs {
			server.Database.FirstOrCreate(&config)
		}
	} else {
		p.Configs = configs
	}

	if err := p.FetchGachaLog(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	// Timestamp in gacha log are all in UTC-8
	timezone, _ := time.LoadLocation("Asia/Shanghai")

	for gachaConfigKey, gachaLogs := range p.GachalLogInPool {
		if len(gachaLogs) == 0 {
			continue
		}
		UID := gachaLogs[0].UID
		var lastWish WishLog
		server.Database.Where(
			map[string]interface{}{"gacha_type": gachaConfigKey, "user_id": UID},
		).Last(&lastWish)

		pityStar4 := 1
		if lastWish.PityStar4 > pityStar4 {
			pityStar4 = lastWish.PityStar4
		}
		pityStar5 := 1
		if lastWish.PityStar5 > pityStar5 {
			pityStar5 = lastWish.PityStar5
		}
		for _, gachaLog := range gachaLogs {
			if gachaLog.ID == lastWish.ID {
				break
			}

			tm, err := time.ParseInLocation("2006-01-02 15:04:05", gachaLog.Time, timezone)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
			}

			var gachaItem parser.GachaItem
			server.Database.Where(&parser.GachaItem{Name: gachaLog.Name}).Last(&gachaItem)
			server.Database.Create(WishLog{
				ID:        gachaLog.ID,
				UserID:    gachaLog.UID,
				GachaType: gachaLog.GachaType,
				ItemID:    gachaItem.ID,
				Time:      tm,
				PityStar4: pityStar4,
				PityStar5: pityStar5,
			})
			// According to study nga bbs. Pity count for 4 star and 5 star items are separated. i.e. this situation can happen 9 3-star items, follows 1 5-star item, then a 4-star item
			if gachaLog.RankType != "5" {
				pityStar5 += 1
			} else {
				pityStar5 = 1
			}

			if gachaLog.RankType != "4" {
				pityStar4 += 1
			} else {
				pityStar4 = 1
			}
		}
	}

}

// GetLogs gets logs from Database
func (server *Server) GetLogs(ctx *gin.Context) {
	UID := ctx.Param("uid")
	rarity := ctx.Query("rarity") // rank_type
	gachaType := ctx.Query("gachaType")
	itemType := ctx.Query("itemType")

	var logs []parser.GachaLog

	result := server.Database.Model(&parser.GachaLog{}).Where(&parser.GachaLog{
		RankType:  rarity,
		UID:       UID,
		GachaType: gachaType,
		ItemType:  itemType,
	}).Find(&logs)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": logs,
	})
}
