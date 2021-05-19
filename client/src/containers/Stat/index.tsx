import React, { useEffect, useState, useRef } from 'react';
import { Row, Col } from 'antd';
import axios from 'axios';
import * as echarts from 'echarts';
import 'echarts/lib/chart/line';
import 'echarts/lib/chart/bar';
import styles from './index.module.less';
import { WishLog } from 'genshin-wish';
import ItemList from '../../components/ItemList';
import StatisticsNumbers from './Statistics';
import { useParams } from 'react-router';

const client = axios.create({
  baseURL: '/api/v1',
});

const Stat: React.FC = () => {
  const { userId, configKey } =
    useParams<{ userId: string; configKey: string }>();
  const [wishLogs, setWishLogs] = useState<WishLog[]>([]);
  const chartRef = useRef<HTMLDivElement>(null);
  let chartInstance: echarts.ECharts;
  useEffect(() => {
    const fetchLog = async () => {
      const gachaLog = await client.get<{ data: WishLog[] }>(`log/${userId}`, {
        params: {
          gachaType: configKey === 'all' ? '' : configKey,
        },
      });
      setWishLogs(gachaLog.data.data);
    };
    fetchLog();

    return () => {
      chartInstance && chartInstance.dispose();
    };
  }, []);

  useEffect(() => {
    if (!chartRef.current || wishLogs.length === 0) {
      return;
    }
    let chartInstance = echarts.getInstanceByDom(chartRef.current);
    if (!chartInstance) {
      chartInstance = echarts.init(chartRef.current);
    }

    chartInstance.setOption({
      title: {
        text: '抽卡动态',
      },
      xAxis: {
        type: 'category',
        name: '物品',
        data: wishLogs.map(log => log.Item.name),
      },
      yAxis: {
        name: '抽数',
      },
      series: [{ type: 'bar', data: wishLogs.map(log => log.pityStar5) }],
    });
  });

  // const chartOption

  return (
    <Row>
      <Col span={8}>
        <ItemList />
      </Col>
      <Col span={16}>
        <StatisticsNumbers wishLogs={wishLogs} />
        <div className={styles.pullChart} ref={chartRef} />
      </Col>
    </Row>
  );
};

export default Stat;
