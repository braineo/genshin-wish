#!/usr/bin/env python3
import requests
import argparse
import urllib.parse
from datetime import datetime
import itertools
import time

argParser = argparse.ArgumentParser(description="genshin wish history")
argParser.add_argument(
    "--forceItem",
    action="store",
    type=bool,
    dest="forceItem",
    default=False,
    help="Force updating item list",
)


class GenshinWishParser(object):

    GachaTypes = {
        "200": "常驻池",
        "301": "角色池",
    }

    def __init__(self, authkey: str) -> None:
        self._authkey = authkey
        self._itemEndpoint = "https://webstatic-sea.mihoyo.com/hk4e/gacha_info/os_asia/items/zh-cn.json"
        self._historyEndpoint = (
            "https://hk4e-api-os.mihoyo.com/event/gacha_info/api/getGachaLog"
        )
        self._itemTable = {}  # item id: item
        self._wishList = {}  # gacha type: wishes

    def get_item_table(self):
        timestamp = datetime.timestamp(datetime.now())
        params = {"ts": int(timestamp) // 100}
        response = requests.get(self._itemEndpoint, params=params)
        # expect a list of dict
        # {
        #  "item_id": "1022",
        #  "name": "温迪",
        #  "item_type": "角色",
        #  "rank_type": "5"
        # },
        for item in response.json():
            self._itemTable[item["item_id"]] = item

    def get_wish_list(self):
        if not self._itemTable:
            self.get_item_table()
        for gachaType in GenshinWishParser.GachaTypes.keys():
            self._wishList[gachaType] = self._get_wish_list(gachaType)
            for wish in self._wishList[gachaType]:
                item = self._itemTable[wish["item_id"]]
                wish["name"] = item["name"]
                wish["rank"] = int(item["rank_type"])
                wish["type"] = item["item_type"]

    def _get_wish_list(self, gachaType: str):
        wishList = []
        for pageNumber in itertools.count(1):
            _wishList = self._get_wish_list_at(pageNumber, gachaType)
            time.sleep(0.1)
            if len(_wishList) > 0:
                wishList.extend(_wishList)
            else:
                break
        return wishList

    def _get_wish_list_at(self, pageNumber: int, gachaType: str):
        ext = '{"loc":{"x":1934.4156494140625,"y":196.4535369873047,"z":-1266.369873046875},"platform":"IOS"}'

        params = {
            "authkey_ver": 1,
            "sign_type": 2,
            "auth_appid": "webview_gacha",
            "init_type": 301,
            "gacha_id": "2ffa459718702872a52867fa0521e32b6843b0",
            "lang": "zh-cn",
            "device_type": "mobile",
            "ext": ext,
            "game_version": "OSRELiOS1.2.0_R1771533_S1847412_D1816310",
            "region": "os_asia",
            "authkey": self._authkey,
            "game_biz": "hk4e_global",
            "gacha_type": gachaType,
            "page": pageNumber,
            "size": 20,  # seems support up to 20 items
        }
        headers = {
            "Cookie": "mi18nLang=zh-cn",
            "Host": "hk4e-api-os.mihoyo.com",
            "Origin": "https://webstatic-sea.mihoyo.com",
        }
        response = requests.Session().get(
            self._historyEndpoint,
            params=urllib.parse.urlencode(params, safe=""),
            headers=headers,
        )

        resData = response.json()
        if not resData["data"]:
            print(resData["message"])
            return []
        else:
            return resData.get("data", {}).get("list", [])

    @staticmethod
    def get_rank_statistics(wishList):
        statistics = {
            "total": len(list(wishList)),
        }
        for rank in [3, 4, 5]:
            statistics["%sstar" % rank] = len(
                list(filter(lambda x: x["rank"] == rank, wishList))
            )

        return statistics

    def print_rank_statistics(result):
        print("总抽数: %d" % result["total"])
        for rank in [3, 4, 5]:
            print(
                "%s星数: %d, 占: %2f"
                % (
                    rank,
                    result["%sstar" % rank],
                    result["%sstar" % rank] / result["total"] * 100,
                )
            )

    @staticmethod
    def get_prediction(wishList: list, rank: int, expectation: int):
        noLuckCount = 0
        for wish in wishList:
            if wish["rank"] != rank:
                noLuckCount += 1
            else:
                break

        return {
            "gaChaToGo": expectation - noLuckCount,
            "rank": rank,
            "noLuckCount": noLuckCount,
            "expectation": expectation,
        }

    @staticmethod
    def print_prediction(result):
        print(
            "{rank}星物品已垫{noLuckCount}, 估计还要{gaChaToGo}({expectation})".format(**result)
        )

    @staticmethod
    def get_item_statistics(wishList):
        itemTable = {}
        for wish in wishList:
            if wish["name"] in itemTable:
                itemTable[wish["name"]] += 1
            else:
                itemTable[wish["name"]] = 1
        return itemTable

    def print_statistics(self):

        print("物品统计")

        itemStatistics = GenshinWishParser.get_item_statistics(
            list(itertools.chain(*self._wishList.values()))
        )
        itemRankSortedList = sorted(
            self._itemTable.values(), key=lambda x: -int(x["rank_type"])
        )
        for item in itemRankSortedList:
            if item["name"] in itemStatistics:
                print("%s: %d" % (item["name"], itemStatistics[item["name"]]))

        for gachaType, gachaName in GenshinWishParser.GachaTypes.items():
            print(gachaName)
            wishList = self._wishList[gachaType]
            GenshinWishParser.print_rank_statistics(
                GenshinWishParser.get_rank_statistics(wishList)
            )
            GenshinWishParser.print_prediction(
                GenshinWishParser.get_prediction(wishList, 4, 10)
            )
            GenshinWishParser.print_prediction(
                GenshinWishParser.get_prediction(wishList, 5, 77)
            )

        print("总计")
        GenshinWishParser.print_rank_statistics(
            GenshinWishParser.get_rank_statistics(
                list(itertools.chain(*self._wishList.values()))
            )
        )


if __name__ == "__main__":
    parse = GenshinWishParser(
        "test"
    )
    parse.get_wish_list()
    parse.print_statistics()
