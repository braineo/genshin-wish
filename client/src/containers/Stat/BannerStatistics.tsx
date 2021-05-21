import React, { useEffect, useState, useRef } from 'react';
import { Row, Col } from 'antd';
import axios from 'axios';
import * as echarts from 'echarts';
import 'echarts/lib/chart/line';
import 'echarts/lib/chart/bar';
import styles from './index.module.less';
import { WishLog } from 'genshin-wish';
import StatisticsNumbers from './StatisticsNumbers';
import Recents from '../Recents';
import { useParams } from 'react-router';

const client = axios.create({
  baseURL: '/api/v1',
});

const BannerStatistics: React.FC = () => {
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
    type Indexable = {
      [key: string]: number;
    };
    const rarityItemCount: Indexable = { '3': 0, '4': 0, '5': 0 };

    chartInstance.setOption({
      title: {
        text: '物品分布',
      },

      series: [
        {
          type: 'pie',
          radius: [20, 60],
          left: 'center',
          height: '100%',
          width: '100%',
          label: {
            alignTo: 'edge',
            formatter: '{name|{b}}\n{time|{c} %}',
            minMargin: 5,
            edgeDistance: 10,
            lineHeight: 15,
            rich: {
              time: {
                fontSize: 10,
                color: '#999',
              },
            },
          },
          data: Object.entries(
            wishLogs.reduce((prev, current) => {
              prev[current.Item.rarity] = prev[current.Item.rarity] + 1;
              return prev;
            }, rarityItemCount),
          ).map(elem => ({
            name: `${elem[0]}星物品`,
            value: Number(((elem[1] / wishLogs.length) * 100).toPrecision(3)),
          })),
        },
      ],
    });
  });

  // const chartOption

  return (
    <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }}>
      <Col xs={16} sm={12} md={8} lg={8} xl={6}>
        <Recents />
      </Col>
      <Col xs={8} sm={12} md={16} lg={16} xl={18}>
        <StatisticsNumbers wishLogs={wishLogs} />
        <div className={styles.pullChart} ref={chartRef} />
      </Col>
    </Row>
  );
};

export default BannerStatistics;
