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
		v1Route.POST("/log", server.PostLogHandler)
		v1Route.GET("/log", GetLog)
	}
	return server
}

func (server *Server) Run() {

	server.Engine.Run(":8080")
}

// PostLogHandler accept query URL for gacha log to query game server
func (server *Server) PostLogHandler(ctx *gin.Context) {
	rawQuery := ctx.PostForm("query")
	playerName := ctx.PostForm("playerName")

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
		var lastLog GachaLog
		server.Database.Where(
			map[string]interface{}{"gacha_type": gachaConfigKey, "player_name": playerName},
		).Last(&lastLog)
		for _, gachaLog := range gachaLogs {
			if gachaLog.Time == lastLog.Time {
				break
			}
			server.Database.Create(GachaLog{
				GachaLog:   gachaLog,
				PlayerName: playerName,
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"test":   rawQuery,
		"player": playerName,
	})

}

// GetLog gets logs from Database
func GetLog(ctx *gin.Context) {

}
