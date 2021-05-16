package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/op/go-logging"

	"github.com/braineo/genshin-wish/parser"
)

var log = logging.MustGetLogger("cli")

func ExportGachaLog(p *parser.GenshinWishParser, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	log.Infof("正在保存到%s ...", filePath)

	err = file.Truncate(0)
	if err != nil {
		return nil
	}

	writer := csv.NewWriter(file)
	writer.Write([]string{"UID", "物品ID", "物品名", "稀有度", "卡池", "抽卡时间"})
	for _, gachaConfig := range p.Configs {
		gachaLogs := p.GachalLogInPool[gachaConfig.Key]
		for _, gachaLog := range gachaLogs {
			itemInfo := p.ItemTable[gachaLog.ItemID]
			writer.Write([]string{gachaLog.UID, itemInfo.ID, itemInfo.Name, itemInfo.RankType, gachaConfig.Name, gachaLog.Time})
		}
	}

	writer.Flush()
	return nil
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

	parser, err := parser.New(args[1])
	if err != nil {
		panic(err)
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
	// export log to file name in current path
	if len(args) == 3 {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		fmt.Println(exPath)
		ExportGachaLog(parser, path.Join(exPath, args[2]))
	}
}
