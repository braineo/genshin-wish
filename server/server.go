package server

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/braineo/genshin-wish/parser"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"gorm.io/gorm"
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
		v1Route.GET("/user", server.GetUsers)
		v1Route.PUT("/user/", server.UpdateUser)

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

func (server *Server) GetUsers(ctx *gin.Context) {

	var users []User
	server.Database.Model(&User{}).Find(&users)

	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func (server *Server) UpdateUser(ctx *gin.Context) {
	var updatedUser User
	if ctx.ShouldBind(&updatedUser) == nil {
		if updatedUser.ID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "ID cannot be empty",
			})
			return
		}
		var user User
		server.Database.First(&user, "id = ?", updatedUser.ID)
		if user.ID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Cannot find user %s", updatedUser.ID),
			})
			return
		}
		server.Database.Model(&user).Updates(updatedUser)
	}
}

func (server *Server) FetchGachaItems(ctx *gin.Context) {
	rawQuery := ctx.PostForm("query")
	p, err := parser.New(rawQuery, parser.WithLanguage(parser.EnUs))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
	var configs []parser.GachaConfig
	server.Database.Model(&parser.GachaConfig{}).Find(&configs)
	ctx.JSON(http.StatusOK, gin.H{
		"data": configs,
	})
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

	for _, gachaLogs := range p.GachalLogInPool {
		if err := server.createWishLogs(&gachaLogs); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
	}
}

func (server *Server) createWishLogs(gachaLogs *[]parser.GachaLog) error {
	if len(*gachaLogs) == 0 {
		return nil
	}

	// Timestamp in gacha log are all in UTC-8
	timezone, _ := time.LoadLocation("Asia/Shanghai")

	UID := (*gachaLogs)[0].UID
	gachaConfigKey := (*gachaLogs)[0].GachaType
	var lastWish WishLog

	endIndex := -1
	pityStar4 := 1
	pityStar5 := 1

	// determine range of new values
	if err := server.Database.Where(
		map[string]interface{}{"gacha_type": gachaConfigKey, "user_id": UID},
	).Last(&lastWish); err != nil {
		log.Debug("no record found, use all logs")
		endIndex = len(*gachaLogs) - 1
	} else {
		log.Debugf("last saved wish ID %s", lastWish.ID)
		pityStar4 = lastWish.PityStar4
		pityStar5 = lastWish.PityStar5

		for _, gachaLog := range *gachaLogs {
			if gachaLog.ID == lastWish.ID {
				break
			}
			endIndex += 1
		}
	}
	reg, _ := regexp.Compile("[^a-zA-Z]+")
	// gacha log parsed in time desc order, process from backwards
	for index := endIndex; index >= 0; index-- {
		gachaLog := (*gachaLogs)[index]

		tm, err := time.ParseInLocation("2006-01-02 15:04:05", gachaLog.Time, timezone)
		if err != nil {
			return err
		}
		itemId := strings.ToLower(reg.ReplaceAllString(gachaLog.Name, ""))
		server.Database.FirstOrCreate(&WishLog{
			ID:     gachaLog.ID,
			UserID: gachaLog.UID,
			User: User{
				ID: gachaLog.UID,
			},
			GachaType: gachaLog.GachaType,
			ItemID:    itemId,
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
	return nil
}

// GetLogs gets logs from Database
func (server *Server) GetLogs(ctx *gin.Context) {
	UID := ctx.Param("uid")
	rarity := ctx.Query("rarity") // rank_type
	gachaType := ctx.Query("gachaType")
	itemType := ctx.Query("itemType")
	size := ctx.Query("size")
	orderBy := ctx.Query("orderBy")
	sortOrder := ctx.Query("sort")

	// cursor := ctx.Query("cursor")

	var logs []WishLog
	// inner join, Item is struct field not the type
	db := server.Database.Debug().Joins("Item").Where(
		&WishLog{
			UserID:    UID,
			GachaType: gachaType,
		})
	querySplitFn := func(c rune) bool {
		return c == '+'
	}
	rarities := strings.FieldsFunc(rarity, querySplitFn)
	if len(rarities) == 1 {
		db = db.Where("item__rarity = ?", rarities[0])
	} else if len(rarities) > 1 {
		db = db.Where("item__rarity in (?)", rarities)
	}

	if size != "" {
		limit, err := strconv.Atoi(size)
		if err == nil {
			db = db.Limit(limit)
		}
	}
	if itemType != "" {
		db = db.Where("Item__type = ?", itemType)
	}

	if orderBy == "" {
		orderBy = "id"
	}
	orderBy = fmt.Sprintf("wish_logs.%s", orderBy)
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	db = db.Order(strings.Join([]string{orderBy, sortOrder}, " "))

	result := db.Find(&logs)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": logs,
	})
}
