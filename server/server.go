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

	v1Route := engine.Group("api/v1")

	{
		v1Route.GET("/user", server.GetUsers)
		v1Route.PUT("/user/", server.UpdateUser)
		// Wish log related
		v1Route.POST("/log", server.FetchLogs)
		v1Route.GET("/log/:uid", server.GetLogs)
		v1Route.GET("/stat/:uid", server.GetStat)

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

func (server *Server) GetGachaConfigs(ctx *gin.Context) {
	var configs []parser.GachaConfig
	server.Database.Model(&parser.GachaConfig{}).Find(&configs)
	ctx.JSON(http.StatusOK, gin.H{
		"data": configs,
	})
}

type queryInfo struct {
	Query string `json:"query"`
}

// FetchLogs accept query URL for gacha log to query game server
func (server *Server) FetchLogs(ctx *gin.Context) {
	var query queryInfo
	if err := ctx.ShouldBind(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	p, err := parser.New(query.Query, parser.WithLanguage(parser.EnUs))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	userID, err := p.GetUserID()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	// all types of gachas. need to specify manually, , 301(角色, with 400),
	gachaTypes := map[string][]string{
		"200": {"200"},        // 200(奔行)
		"301": {"301", "400"}, // 301(角色) 400(角色2)
		"302": {"302"},        // 302(武器)
	}

	for gachaType, queryGachaTypes := range gachaTypes {
		endId := ""
		var lastWish WishLog
		if err := server.Database.Where(
			map[string]interface{}{"gacha_type": queryGachaTypes, "user_id": userID},
		).Last(&lastWish); err != nil {
			log.Debugf("found last record 5 star pity %d, ID %s", lastWish.PityStar5, lastWish.ID)
			endId = lastWish.ID
			ctx.JSON(http.StatusOK, gin.H{
				"wish": lastWish,
			})
			return
		}

		if gachaLogs, err := p.FetchGachaLog(gachaType, endId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})

			if err := server.createWishLogs(gachaLogs, gachaType); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
			}
		}
	}
}

func (server *Server) createWishLogs(gachaLogs *[]parser.GachaLog, gachaType string) error {
	if len(*gachaLogs) == 0 {
		return nil
	}

	// Timestamp in gacha log are all in UTC-8
	timezone, _ := time.LoadLocation("Asia/Shanghai")

	UID := (*gachaLogs)[0].UID
	var lastWish WishLog

	endIndex := -1
	pityStar4 := 1
	pityStar5 := 1

	// determine range of new values
	if err := server.Database.Where(
		map[string]interface{}{"gacha_type": gachaType, "user_id": UID},
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
			endIndex++
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
			pityStar5++
		} else {
			pityStar5 = 1
		}

		if gachaLog.RankType != "4" {
			pityStar4++
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

	var logs []WishLog
	// inner join, Item is struct field not the type
	db := server.Database.Joins("Item").Where(
		&WishLog{
			UserID: UID,
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

	gachaTypes := strings.FieldsFunc(gachaType, querySplitFn)
	if len(gachaTypes) == 1 {
		db = db.Where("gacha_type = ?", gachaTypes[0])
	} else if len(gachaTypes) > 1 {
		db = db.Where("gacha_type in (?)", gachaTypes)
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

// GetStat gets statisitcs of wishes, including repeat of items
func (server *Server) GetStat(ctx *gin.Context) {
	UID := ctx.Param("uid")
	gachaType := ctx.Query("gachaType")
	itemType := ctx.Query("itemType")

	// inner join, Item is struct field not the type
	db := server.Database.Joins("Item").Where(
		&WishLog{
			UserID:    UID,
			GachaType: gachaType,
		})

	if itemType != "" {
		db = db.Where("Item__type = ?", itemType)
	}

	var logs []ItemCount

	db.Model(&WishLog{}).Select("*, count(item_id) as count").Group("item_id").Find(&logs)

	ctx.JSON(http.StatusOK, gin.H{
		"data": logs,
	})
}
