package parser

import (
	"fmt"
	"math"
	"sort"
)

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
	Star5Intervals        []int
	ShortestStar5Interval int
	LongestStar5Interval  int
	CurrentStar5Interval  int
	CurrentStar4Interval  int
	ItemCount             map[string]int
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
			ShortestStar5Interval: 90,
			LongestStar5Interval:  0,
			CurrentStar5Interval:  0,
			CurrentStar4Interval:  0,
			ItemCount:             make(map[string]int),
		}

		for _, gachaLog := range gachaLogs {
			statistics.Total++
			p.Statistics.ItemCount[gachaLog.Name]++

			isCharacter := true
			if gachaLog.ItemType == "角色" {
				statistics.Character++
			} else if gachaLog.ItemType == "武器" {
				statistics.Weapon++
				isCharacter = false
			}

			if gachaLog.RankType == "5" {
				statistics.Star5++
				if isCharacter {
					statistics.CharacterStar5++
				} else {
					statistics.WeaponStar5++
				}
				if foundFirstStar5Item {
					statistics.Star5Intervals = append(statistics.Star5Intervals, star5Interval)
					statistics.LongestStar5Interval = int(math.Max(float64(star5Interval), float64(statistics.LongestStar5Interval)))
					statistics.ShortestStar5Interval = int(math.Min(float64(star5Interval), float64(statistics.ShortestStar5Interval)))
				}
				foundFirstStar5Item = true
				star5Interval = 0
			} else if gachaLog.RankType == "4" {
				statistics.Star4++
				if isCharacter {
					statistics.CharacterStar4++
				} else {
					statistics.WeaponStar4++
				}
				foundFirstStar4Item = true
				star5Interval++
			} else if gachaLog.RankType == "3" {
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

		if foundFirstStar5Item {
			statistics.Star5Intervals = append(statistics.Star5Intervals, star5Interval)
			statistics.LongestStar5Interval = int(math.Max(float64(star5Interval), float64(statistics.LongestStar5Interval)))
			statistics.ShortestStar5Interval = int(math.Min(float64(star5Interval), float64(statistics.ShortestStar5Interval)))
		}
		p.StatisticsInPool[gachaKey] = statistics

		p.Statistics.Total += statistics.Total
		p.Statistics.Star5 += statistics.Star5
		p.Statistics.Star4 += statistics.Star4
		p.Statistics.Star3 += statistics.Star3
		p.Statistics.Character += statistics.Character
		p.Statistics.CharacterStar5 += statistics.CharacterStar5
		p.Statistics.CharacterStar4 += statistics.CharacterStar4
		p.Statistics.Weapon += statistics.Weapon
		p.Statistics.WeaponStar5 += statistics.WeaponStar5
		p.Statistics.WeaponStar4 += statistics.WeaponStar4
		p.Statistics.LongestStar5Interval = int(math.Max(float64(p.Statistics.LongestStar5Interval), float64(statistics.LongestStar5Interval)))
		p.Statistics.ShortestStar5Interval = int(math.Min(float64(p.Statistics.ShortestStar5Interval), float64(statistics.ShortestStar5Interval)))
		p.Statistics.Star5Intervals = append(p.Statistics.Star5Intervals, statistics.Star5Intervals...)
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
			statistics.WeaponStar4,
			float32(statistics.WeaponStar4)/float32(statistics.Weapon)*100.0,
		)
		fmt.Printf("四星物品已垫%d,估计还要%d(%d)\n", statistics.CurrentStar4Interval, 10-statistics.CurrentStar4Interval, 10)
		fmt.Printf("五星物品已垫%d,估计还要%d(%d)\n", statistics.CurrentStar5Interval, 77-statistics.CurrentStar5Interval, 77)
		if statistics.ShortestStar5Interval < 90 {
			fmt.Printf("最短五星抽数%d,最长五星抽数%d,平均%.2f\n", statistics.ShortestStar5Interval, statistics.LongestStar5Interval, mean(statistics.Star5Intervals))
		}
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
		p.Statistics.WeaponStar4,
		float32(p.Statistics.WeaponStar4)/float32(p.Statistics.Total)*100.0,
	)
	fmt.Printf("最短五星抽数%d,最长五星抽数%d,平均%.2f\n", p.Statistics.ShortestStar5Interval, p.Statistics.LongestStar5Interval, mean(p.Statistics.Star5Intervals))

	fmt.Println("==========")
	fmt.Println("物品统计")

	itemSlice := make([]GachaItem, 0, len(p.ItemTable))

	for _, item := range p.ItemTable {
		itemSlice = append(itemSlice, item)
	}
	sort.Slice(itemSlice, func(i, j int) bool {
		if itemSlice[i].ItemType < itemSlice[j].ItemType {
			return true
		}
		if itemSlice[i].ItemType > itemSlice[j].ItemType {
			return false
		}
		return itemSlice[i].RankType > itemSlice[j].RankType
	})

	for _, item := range itemSlice {
		if p.Statistics.ItemCount[item.Name] > 0 {
			fmt.Printf("%s: %d\n", item.Name, p.Statistics.ItemCount[item.Name])
		}
	}
}

func mean(intSlice []int) float32 {
	if len(intSlice) == 0 {
		return 0
	}
	sum := 0
	for _, number := range intSlice {
		sum += number
	}

	return float32(sum) / float32(len(intSlice))
}
