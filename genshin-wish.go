package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

type GachaItem struct {
	ItemID   string `json:"item_id"`
	Name     string `json:"name"`
	ItemType string `json:"item_type"`
	RankType string `json:"rank_type"`
}

type GachaConfig struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type GachaConfigResponse struct {
	RetCode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		GachaTypeList []GachaConfig `json:"gacha_type_list"`
	} `json:"data"`
}

type GachaLog struct {
	GachaType string `json:"gacha_type"`
	UID       string `json:"uid"`
	Count     string `json:"conut"`
	ItemID    string `json:"item_id"`
	Time      string `json:"time"`
}

type GachaLogResponse struct {
	RetCode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		GachaLogList []GachaLog `json:"list"`
	} `json:"data"`
}

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

type GachaStatistics struct {
	Total                 int
	Star5                 int
	Star4                 int
	Star3                 int
	Character             int
	CharacterStar5        int
	CharacterStar4        int
	Weapon                int
	WeaponStar5           int
	WeaponStar4           int
	ShortestStar5Interval int
	LongestStar5Interval  int
	CurrentStar5Interval  int
	CurrentStar4Interval  int
	ItemCount             map[string]int
}

const (
	itemListURL    = "https://webstatic-sea.mihoyo.com/hk4e/gacha_info/os_asia/items/zh-cn.json"
	gachaConfigURL = "https://hk4e-api-os.mihoyo.com/event/gacha_info/api/getConfigList"
	gachaLogURL    = "https://hk4e-api-os.mihoyo.com/event/gacha_info/api/getGachaLog"
)

var (
	requiredQueryFields = []string{"authkey_ver", "sign_type", "auth_appid", "gacha_id", "lang", "game_biz", "authkey", "region"}
)

func (p *GenshinWishParser) prepareRequestParams(request *http.Request) url.Values {
	query := request.URL.Query()
	for _, name := range requiredQueryFields {
		query.Set(name, p.Query[name][0])
	}
	return query

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

func (p *GenshinWishParser) MakeStatistics() {

	for gachaKey, gachaLogs := range p.GachalLogInPool {
		foundFirstStar5Item := false
		foundFirstStar4Item := false

		star5Interval := 0

		statistics := GachaStatistics{
			Total:                 0,
			Star5:                 0,
			Star4:                 0,
			Star3:                 0,
			Character:             0,
			CharacterStar5:        0,
			CharacterStar4:        0,
			Weapon:                0,
			WeaponStar5:           0,
			WeaponStar4:           0,
			ShortestStar5Interval: 0,
			LongestStar5Interval:  0,
			CurrentStar5Interval:  0,
			CurrentStar4Interval:  0,
			ItemCount:             make(map[string]int),
		}

		for _, gachaLog := range gachaLogs {
			itemInfo := p.ItemTable[gachaLog.ItemID]
			statistics.Total++
			p.Statistics.ItemCount[itemInfo.Name]++

			isCharacter := true
			if itemInfo.ItemType == "角色" {
				statistics.Character++
			} else if itemInfo.ItemType == "武器" {
				statistics.Weapon++
				isCharacter = false
			}

			if itemInfo.RankType == "5" {
				statistics.Star5++
				if isCharacter {
					statistics.CharacterStar5++
				} else {
					statistics.WeaponStar5++
				}
				foundFirstStar5Item = true
				statistics.LongestStar5Interval = int(math.Max(float64(star5Interval), float64(statistics.LongestStar5Interval)))
				statistics.ShortestStar5Interval = int(math.Min(float64(star5Interval), float64(statistics.ShortestStar5Interval)))
				star5Interval = 0
			} else if itemInfo.RankType == "4" {
				statistics.Star4++
				if isCharacter {
					statistics.CharacterStar4++
				} else {
					statistics.WeaponStar4++
				}
				foundFirstStar4Item = true
				star5Interval++
			} else if itemInfo.RankType == "3" {
				statistics.Star3++
				star5Interval++
			}

			if !foundFirstStar5Item {
				statistics.CurrentStar5Interval++

			}
			if !foundFirstStar4Item {
				statistics.CurrentStar4Interval++
			}
		}
		p.StatisticsInPool[gachaKey] = statistics

		p.Statistics.Total = p.Statistics.Total + statistics.Total
		p.Statistics.Star5 = p.Statistics.Star5 + statistics.Star5
		p.Statistics.Star4 = p.Statistics.Star4 + statistics.Star4
		p.Statistics.Star3 = p.Statistics.Star3 + statistics.Star3
		p.Statistics.Character = p.Statistics.Character + statistics.Character
		p.Statistics.CharacterStar5 = p.Statistics.CharacterStar5 + statistics.CharacterStar5
		p.Statistics.CharacterStar4 = p.Statistics.Character + statistics.Character
		p.Statistics.Weapon = p.Statistics.Weapon + statistics.Weapon
		p.Statistics.WeaponStar5 = p.Statistics.WeaponStar5 + statistics.WeaponStar5
		p.Statistics.WeaponStar4 = p.Statistics.WeaponStar4 + statistics.WeaponStar4
	}
}

func (p *GenshinWishParser) PrintStatistics() {
	for _, gachaConfig := range p.Configs {
		fmt.Println("==========")
		fmt.Printf("%s抽卡统计\n", gachaConfig.Name)
		statistics := p.StatisticsInPool[gachaConfig.Key]
		if statistics.Total == 0 {
			continue
		}
		fmt.Printf("总数%v 五星%v(%.2f%%) 四星%v(%.2f%%) 三星%v(%.2f%%)\n",
			statistics.Total,
			statistics.Star5,
			float32(statistics.Star5)/float32(statistics.Total)*100.0,
			statistics.Star4,
			float32(statistics.Star4)/float32(statistics.Total)*100.0,
			statistics.Star3,
			float32(statistics.Star3)/float32(statistics.Total)*100.0,
		)
		fmt.Printf("角色%v 五星%v(%.2f%%) 四星%v(%.2f%%)\n",
			statistics.Character,
			statistics.CharacterStar5,
			float32(statistics.CharacterStar5)/float32(statistics.Character)*100.0,
			statistics.CharacterStar4,
			float32(statistics.CharacterStar4)/float32(statistics.Character)*100.0,
		)
		fmt.Printf("武器%v 五星%v(%.2f%%) 四星%v(%.2f%%)\n",
			statistics.Weapon,
			statistics.WeaponStar5,
			float32(statistics.WeaponStar5)/float32(statistics.Weapon)*100.0,
			statistics.CharacterStar4,
			float32(statistics.WeaponStar4)/float32(statistics.Weapon)*100.0,
		)
		fmt.Printf("四星物品已垫%d,估计还要%d(%d)\n", statistics.CurrentStar4Interval, 10-statistics.CurrentStar4Interval, 10)
		fmt.Printf("五星物品已垫%d,估计还要%d(%d)\n", statistics.CurrentStar5Interval, 77-statistics.CurrentStar5Interval, 77)
	}
	fmt.Println("==========")
	fmt.Println("综合统计")
	fmt.Printf("总数%v 五星%v(%.2f%%) 四星%v(%.2f%%) 三星%v(%.2f%%)\n",
		p.Statistics.Total,
		p.Statistics.Star5,
		float32(p.Statistics.Star5)/float32(p.Statistics.Total)*100.0,
		p.Statistics.Star4,
		float32(p.Statistics.Star4)/float32(p.Statistics.Total)*100.0,
		p.Statistics.Star3,
		float32(p.Statistics.Star3)/float32(p.Statistics.Total)*100.0,
	)
	fmt.Printf("角色%v(%.2f%%) 五星%v(%.2f%%) 四星%v(%.2f%%)\n",
		p.Statistics.Character,
		float32(p.Statistics.Character)/float32(p.Statistics.Total)*100.0,
		p.Statistics.CharacterStar5,
		float32(p.Statistics.CharacterStar5)/float32(p.Statistics.Total)*100.0,
		p.Statistics.CharacterStar4,
		float32(p.Statistics.CharacterStar4)/float32(p.Statistics.Total)*100.0,
	)
	fmt.Printf("武器%v(%.2f%%) 五星%v(%.2f%%) 四星%v(%.2f%%)\n",
		p.Statistics.Weapon,
		float32(p.Statistics.Weapon)/float32(p.Statistics.Total)*100.0,
		p.Statistics.WeaponStar5,
		float32(p.Statistics.WeaponStar5)/float32(p.Statistics.Total)*100.0,
		p.Statistics.CharacterStar4,
		float32(p.Statistics.WeaponStar4)/float32(p.Statistics.Total)*100.0,
	)
	fmt.Println("==========")
	fmt.Println("物品统计")

	itemSlice := make([]GachaItem, 0, len(p.ItemTable))

	for _, item := range p.ItemTable {
		itemSlice = append(itemSlice, item)
	}
	sort.Slice(itemSlice, func(i, j int) bool {
		return itemSlice[i].RankType > itemSlice[j].RankType
	})

	for _, item := range itemSlice {
		if p.Statistics.ItemCount[item.Name] > 0 {
			fmt.Printf("%s: %d\n", item.Name, p.Statistics.ItemCount[item.Name])
		}
	}
}

func main() {
	formatter := logging.MustStringFormatter("%{color}%{time:2006-01-02T15:04:05.000000-07:00} %{module} [%{level}] <%{pid}> [%{shortfile} %{shortfunc}] %{message}%{color:reset}")
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), formatter))
	if level, err := logging.LogLevel("Info"); err == nil {
		logging.SetLevel(level, "")
	}
	args := os.Args
	if len(args) < 2 {
		log.Error("需要authkey链接")
	}

	u, err := url.Parse(args[1])

	if err != nil {
		panic(err)
	}

	query, _ := url.ParseQuery(u.RawQuery)

	parser := GenshinWishParser{
		Client:           http.Client{},
		Query:            query,
		ItemTable:        make(map[string]GachaItem),
		Configs:          make([]GachaConfig, 0),
		GachalLogInPool:  make(map[string][]GachaLog),
		StatisticsInPool: make(map[string]GachaStatistics),
		Statistics:       GachaStatistics{ItemCount: make(map[string]int)},
	}

	err = parser.FetchGachaConfigs()
	if err != nil {
		panic(err)
	}
	log.Debug(parser.Configs)
	err = parser.FetchGachaItems()
	if err != nil {
		panic(err)
	}
	log.Debug(parser.ItemTable)

	err = parser.FetchGachaLog()
	if err != nil {
		panic(err)
	}
	parser.MakeStatistics()
	log.Debug(parser.StatisticsInPool)
	parser.PrintStatistics()
}
