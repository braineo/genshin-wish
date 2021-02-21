package parser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/op/go-logging"
)

var (
	log                 = logging.MustGetLogger("parser")
	requiredQueryFields = []string{"authkey_ver", "sign_type", "auth_appid", "lang", "game_biz", "authkey", "region"}
)

const (
	itemListURL    = "https://webstatic-sea.mihoyo.com/hk4e/gacha_info/os_asia/items/zh-cn.json"
	gachaConfigURL = "https://hk4e-api-os.mihoyo.com/event/gacha_info/api/getConfigList"
	gachaLogURL    = "https://hk4e-api-os.mihoyo.com/event/gacha_info/api/getGachaLog"
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
}

func New(rawQuery string) (*GenshinWishParser, error) {

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
	}

	return &parser, nil
}

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

func (p *GenshinWishParser) FetchGachaItems() error {
	log.Infof("正在获取所有物品信息")

	request, err := http.NewRequest("GET", itemListURL, nil)
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
		p.ItemTable[item.ItemID] = item
	}

	return nil
}

func (p *GenshinWishParser) FetchGachaLog() error {
	for _, config := range p.Configs {
		log.Infof("正在获取%s信息", config.Name)
		gachaLog := make([]GachaLog, 0)
		for pageNumber := 1; ; pageNumber++ {
			pagedGachaLog, err := p.fetchGachaLog(pageNumber, config.Key)
			if err != nil {
				log.Debugf("无法读取%s页信息,错误%s", pageNumber, err)
				return err
			}
			if len(pagedGachaLog) == 0 {
				break
			}
			gachaLog = append(gachaLog, pagedGachaLog...)
		}
		p.GachalLogInPool[config.Key] = gachaLog
	}
	return nil
}

func (p *GenshinWishParser) fetchGachaLog(pageNumber int, gachaType string) ([]GachaLog, error) {
	request, err := http.NewRequest("GET", gachaLogURL, nil)
	if err != nil {
		return nil, err
	}
	query := p.prepareRequestParams(request)
	query.Set("page", strconv.Itoa(pageNumber))
	query.Set("size", "20")
	query.Set("gacha_type", gachaType)
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
	return query

}
