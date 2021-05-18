import React, { useEffect, useState, useRef } from "react";
import axios from "axios";
import * as echarts from "echarts";
import "echarts/lib/chart/line";
import "echarts/lib/chart/bar";
import styles from "./App.module.scss";

type GachaItem = {
  name: string;
  type: "weapon" | "character";
  rarity: string;
};
type GachaLog = {
  gachaType: string;
  id: string;
  Item: GachaItem;
  time: string;
  pityStar4: number;
  pityStar5: number;
};

const client = axios.create({
  baseURL: "/api/v1",
});

function App() {
  const [gachaLogs, setGachaLogs] = useState<GachaLog[]>([]);
  const chartRef = useRef<HTMLDivElement>(null);
  let chartInstance: echarts.ECharts;
  useEffect(() => {
    const fetchLog = async () => {
      const gachaLog = await client.get<{ data: GachaLog[] }>("log/820575774", {
        params: {
          rarity: "5",
          itemType: "character",
        },
      });
      setGachaLogs(gachaLog.data.data);
    };
    fetchLog();

    return () => {
      chartInstance && chartInstance.dispose();
    };
  }, []);

  useEffect(() => {
    if (!chartRef.current || gachaLogs.length === 0) {
      return;
    }
    let chartInstance = echarts.getInstanceByDom(chartRef.current);
    if (!chartInstance) {
      chartInstance = echarts.init(chartRef.current);
    }

    chartInstance.setOption({
      title: {
        text: "抽卡动态",
      },
      xAxis: {
        type: "category",
        name: "物品",
        data: gachaLogs.map((log) => log.Item.name),
      },
      yAxis: {
        name: "抽数",
      },
      series: [{ type: "bar", data: gachaLogs.map((log) => log.pityStar5) }],
    });
  });

  // const chartOption

  return <div className={styles.pullChart} ref={chartRef} />;
}

export default App;
