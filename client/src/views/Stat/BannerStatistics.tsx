import { Col, Row } from 'antd';
import axios from 'axios';
import * as echarts from 'echarts';
import 'echarts/lib/chart/line';
import 'echarts/lib/chart/bar';
import { WishLog } from 'genshin-wish';
import React, { useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router';
import Recents from '../Recents';
import styles from './index.module.less';
import StatisticsNumbers from './StatisticsNumbers';

const client = axios.create({
  baseURL: '/api/v1',
});

const BannerStatistics: React.FC = () => {
  const { userId, gachaType } = useParams<{
    userId: string;
    gachaType: string;
  }>();
  const [wishLogs, setWishLogs] = useState<WishLog[]>([]);
  const chartRef = useRef<HTMLDivElement>(null);
  let chartInstance: echarts.ECharts;
  useEffect(() => {
    const fetchLog = async () => {
      const gachaLog = await client.get<{ data: WishLog[] }>(`log/${userId}`, {
        params: {
          // FIXME: handle mix gacha pool
          gachaType:
            gachaType === 'all'
              ? ''
              : gachaType === '301'
              ? '301+400'
              : gachaType,
        },
      });
      setWishLogs(gachaLog.data.data);
    };
    fetchLog();

    return () => {
      chartInstance.dispose();
    };
  }, []);

  const chartTitles = [
    {
      text: '物品分布',
    },
    {
      subtext: '物品稀有度分布',
      left: '16.67%',
      top: '75%',
      textAlign: 'center',
    },
    {
      subtext: '物品类别分布',
      left: '50%',
      top: '75%',
      textAlign: 'center',
    },
    {
      text: '抽卡动态',
      top: '66.667%',
    },
  ];

  useEffect(() => {
    if (!chartRef.current || wishLogs.length === 0) {
      return;
    }
    let chartInstance = echarts.getInstanceByDom(chartRef.current);
    chartInstance = echarts.init(chartRef.current);
    type Indexable = {
      [key: string]: number;
    };
    const rarityItemCount: Indexable = { '5': 0, '4': 0, '3': 0 };
    const categoryItemCount = {
      character5Star: 0,
      weapon5Star: 0,
      character4Star: 0,
      weapon4Star: 0,
    };
    const repeatedItemCount: Indexable = {};

    chartInstance.setOption({
      title: chartTitles,
      yAxis: {
        type: 'category',
      },
      xAxis: {
        max: 'dataMax',
      },
      series: [
        {
          type: 'pie',
          radius: [20, 60],
          height: '100%',
          left: 0,
          right: '66.6667%',
          label: {
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

        {
          type: 'pie',
          radius: [20, 60],
          height: '100%',
          left: '33.3333%',
          right: '33.3333%',
          label: {
            formatter: '{name|{b}}\n{time|{c}}',
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
              if (current.Item.type === 'weapon') {
                if (current.Item.rarity === '4') {
                  prev.weapon4Star += 1;
                } else if (current.Item.rarity === '5') {
                  prev.weapon5Star += 1;
                }
              } else {
                if (current.Item.rarity === '4') {
                  prev.character4Star += 1;
                } else if (current.Item.rarity === '5') {
                  prev.character5Star += 1;
                }
              }
              return prev;
            }, categoryItemCount),
          ).map(elem => ({
            name: `${elem[0].endsWith('4Star') ? '四星' : '五星'}${
              elem[0].startsWith('character') ? '人物' : '武器'
            }`,
            value: elem[1],
          })),
        },

        {
          type: 'bar',
          top: '66.667%',
          left: 0,
          height: '100%',
          data: Object.entries(
            wishLogs.reduce((prev, current) => {
              if (current.Item.rarity === '3') {
                return prev;
              }
              if (current.Item.name in prev) {
                prev[current.Item.name] += 1;
              } else {
                prev[current.Item.name] = 1;
              }
              return prev;
            }, repeatedItemCount),
          )
            .filter(elem => elem[1] > 1)
            .sort((a, b) => a[1] - b[1])
            .map(elem => [elem[1], elem[0]]),
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
