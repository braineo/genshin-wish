package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/op/go-logging"
)

var (
	log                 = logging.MustGetLogger("parser")
	requiredQueryFields = []string{"authkey_ver", "sign_type", "auth_appid", "lang", "game_biz", "authkey"}
)

const (
	itemListURLt   = "https://webstatic-sea.mihoyo.com/hk4e/gacha_info/os_asia/items/%s.json"
	gachaConfigURL = "https://hk4e-api-os.mihoyo.com/event/gacha_info/api/getConfigList"
	gachaLogURL    = "https://hk4e-api-os.mihoyo.com/event/gacha_info/api/getGachaLog"
)

type Language string

const (
	ZhCn Language = "zh-cn"
	EnUs Language = "en-us"
)

type GenshinWishParser struct {
	Client           http.Client
	Query            url.Values
	Authkey          string
	ItemTable        map[string]GachaItem
	Configs          []GachaConfig
	GachalLogInPool  map[string][]GachaLog
	StatisticsInPool map[string]GachaStatistics
	Statistics       GachaStatistics
	Language         Language
}

type ParserOptions func(*GenshinWishParser)

func WithLanguage(l Language) ParserOptions {
	return func(gwp *GenshinWishParser) {
		gwp.Language = l
	}
}

// New creates parser from query string
func New(rawQuery string, options ...ParserOptions) (*GenshinWishParser, error) {
	log.Info(rawQuery)

	u, err := url.Parse(rawQuery)

	if err != nil {
		return nil, err
	}

	query, _ := url.ParseQuery(u.RawQuery)

	for _, field := range requiredQueryFields {
		if _, present := query[field]; !present {
			log.Errorf("需要field %v,但不在提供的URL中", field)
		}
	}

	parser := GenshinWishParser{
		Client:           http.Client{},
		Query:            query,
		ItemTable:        make(map[string]GachaItem),
		Configs:          make([]GachaConfig, 0),
		GachalLogInPool:  make(map[string][]GachaLog),
		StatisticsInPool: make(map[string]GachaStatistics),
		Statistics: GachaStatistics{
			ItemCount:             make(map[string]int),
			ShortestStar5Interval: 90,
		},
		Language: ZhCn,
	}

	for _, opt := range options {
		opt(&parser)
	}

	return &parser, nil
}

// Fetch all gacha pools. Seems deprecated by MiHoYo
func (p *GenshinWishParser) FetchGachaConfigs() error {
	log.Infof("正在获取所有卡池列表")

	request, err := http.NewRequest("GET", gachaConfigURL, nil)
	if err != nil {
		return err
	}
	query := p.prepareRequestParams(request)
	request.URL.RawQuery = query.Encode()

	response, err := p.Client.Do(request)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var gachaConfigResponse GachaConfigResponse
	err = json.Unmarshal(body, &gachaConfigResponse)
	if err != nil {
		return err
	}
	if gachaConfigResponse.RetCode == -1 {
		return errors.New(gachaConfigResponse.Message)
	}
	p.Configs = gachaConfigResponse.Data.GachaTypeList
	log.Debugf("%s", gachaConfigResponse.Data.GachaTypeList)
	return nil
}

// Fetch all avaialble gacha items. Seems deprecated by MiHoYo
func (p *GenshinWishParser) FetchGachaItems() error {
	log.Infof("正在获取所有物品信息")

	request, err := http.NewRequest("GET", fmt.Sprintf(itemListURLt, p.Language), nil)
	if err != nil {
		return err
	}
	response, err := p.Client.Do(request)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var gachaItems []GachaItem
	err = json.Unmarshal(body, &gachaItems)
	if err != nil {
		return err
	}
	for _, item := range gachaItems {
		p.ItemTable[item.ID] = item
	}

	return nil
}

// TODO: add stopAt and config, do not loop configs here
func (p *GenshinWishParser) FetchAllGachaLog() error {
	for _, config := range p.Configs {
		log.Infof("正在获取%s信息", config.Name)
		gachaLog := make([]GachaLog, 0)
		endId := "0"
		for pageNumber := 1; ; pageNumber++ {
			pagedGachaLog, err := p.fetchGachaLog(pageNumber, config.Key, endId)
			if err != nil {
				log.Debugf("无法读取%s页信息,错误%s", pageNumber, err)
				return err
			}
			if len(pagedGachaLog) == 0 {
				break
			}
			endId = pagedGachaLog[len(pagedGachaLog)-1].ID
			gachaLog = append(gachaLog, pagedGachaLog...)
			time.Sleep(1 * time.Second)
		}
		log.Debugf("Fetched %v items", len(gachaLog))
		p.GachalLogInPool[config.Key] = gachaLog
	}
	return nil
}

func (p *GenshinWishParser) FetchGachaLog(gachaType string, stopAtId string) error {
	gachaLog := make([]GachaLog, 0)
	endId := "0"
	shouldStop := false
	for pageNumber := 1; ; pageNumber++ {
		pagedGachaLog, err := p.fetchGachaLog(pageNumber, gachaType, endId)
		if err != nil {
			log.Debugf("无法读取%s页信息,错误%s", pageNumber, err)
			return err
		}
		if len(pagedGachaLog) == 0 {
			break
		}
		endId = pagedGachaLog[len(pagedGachaLog)-1].ID
		gachaLog = append(gachaLog, pagedGachaLog...)
		time.Sleep(1 * time.Second)

		for _, gachaLog := range pagedGachaLog {
			if gachaLog.ID == stopAtId {
				shouldStop = true
				break
			}
		}
		if shouldStop {
			break
		}
	}
	log.Debugf("Fetched %v items", len(gachaLog))
	p.GachalLogInPool[gachaType] = gachaLog
	return nil
}

// Fetches gacha log per page
func (p *GenshinWishParser) fetchGachaLog(pageNumber int, gachaType string, endID string) ([]GachaLog, error) {
	request, err := http.NewRequest("GET", gachaLogURL, nil)
	if err != nil {
		return nil, err
	}
	query := p.prepareRequestParams(request)
	query.Set("page", strconv.Itoa(pageNumber))
	query.Set("size", "20")
	query.Set("gacha_type", gachaType)
	query.Set("end_id", endID)
	request.URL.RawQuery = query.Encode()

	response, err := p.Client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var gachaLogResponse GachaLogResponse
	err = json.Unmarshal(body, &gachaLogResponse)
	if err != nil {
		return nil, err
	}
	if gachaLogResponse.RetCode == -1 {
		return nil, errors.New(gachaLogResponse.Message)
	}
	return gachaLogResponse.Data.GachaLogList, nil
}

func (p *GenshinWishParser) prepareRequestParams(request *http.Request) url.Values {
	query := request.URL.Query()
	for _, name := range requiredQueryFields {
		query.Set(name, p.Query[name][0])
	}
	if p.Language != "" {
		query.Set("lang", string(p.Language))
	}
	return query

}
