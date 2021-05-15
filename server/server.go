package server

import (
	"net/http"

	"github.com/braineo/genshin-wish/parser"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func handleCors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}

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
		v1Route.POST("/log", server.FetchLog)
		v1Route.GET("/log/:uid", server.GetLog)
	}
	return server
}

func (server *Server) Run() {

	server.Engine.Run(":8080")
}

// FetchLog accept query URL for gacha log to query game server
func (server *Server) FetchLog(ctx *gin.Context) {
	rawQuery := ctx.PostForm("query")

	p, err := parser.New(rawQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	var configs []parser.GachaConfig
	server.Database.Find(&configs)

	if len(configs) == 0 {
		if err = p.FetchGachaConfigs(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
		// TODO insert new gacha configs
		for _, config := range p.Configs {
			server.Database.Create(&config)
		}
	} else {
		p.Configs = configs
	}

	if err := p.FetchGachaLog(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	for gachaConfigKey, gachaLogs := range p.GachalLogInPool {
		if len(gachaLogs) == 0 {
			continue
		}
		UID := gachaLogs[0].UID
		var lastLog parser.GachaLog
		server.Database.Where(
			map[string]interface{}{"gacha_type": gachaConfigKey, "UID": UID},
		).Last(&lastLog)
		for _, gachaLog := range gachaLogs {
			if gachaLog.Time == lastLog.Time {
				break
			}
			server.Database.Create(gachaLog)
		}
	}

}

// GetLog gets logs from Database
func (server *Server) GetLog(ctx *gin.Context) {
	UID := ctx.Param("uid")
	rarity := ctx.Query("rarity") // rank_type
	gachaType := ctx.Query("gachaType")
	itemType := ctx.Query("itemType")

	var logs []parser.GachaLog

	result := server.Database.Model(&parser.GachaLog{}).Where(&parser.GachaLog{
		RankType: rarity,
		UID:      UID,
		GachaType: gachaType,
		ItemType: itemType,
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
