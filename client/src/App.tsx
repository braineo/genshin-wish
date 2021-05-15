import React, { useEffect, useState, useRef } from "react";
import axios from "axios";
import * as echarts from "echarts";
import "./App.scss";

type GachaLog = {
  gacha_type: string;
  id: string;
  item_type: string;
  name: string;
  rank_type: string;
  time: string;
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
      const gachaLog = await client.get<{ data: GachaLog[] }>("log");
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
        type: "time",
        name: "日期",
      },
      yAxis: {
        name: "抽数",
      },
      series: [{ type: "bar", data: gachaLogs.map((log) => log.rank_type) }],
    });
  });

  // const chartOption

  return <div className="pull-chart" ref={chartRef} />;
}

export default App;
