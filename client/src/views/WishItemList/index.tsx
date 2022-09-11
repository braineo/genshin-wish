import { Table, PageHeader } from 'antd';
import axios from 'axios';
import {
  Character,
  QueryOptions,
  Weapon,
  characters,
  weapons,
} from 'genshin-db';
import { WishLog } from 'genshin-wish';
import { useEffect, useState } from 'react';
import { useHistory, useParams } from 'react-router';
import { Avatar } from '../../components/ItemCard';

type GachaItem = Character | Weapon;

type GachaConfig = {
  name: string;
  key: string;
};

const queryOption = {
  matchAliases: false,
  matchCategories: false,
  verboseCategories: false,
  queryLanguages: ['English'],
  resultLanguage: 'CHS',
};

const client = axios.create({
  baseURL: '/api/v1',
});

export const WishItemList = () => {
  const history = useHistory();
  const { userId } = useParams<{
    userId: string;
  }>();

  return (
    <>
      <PageHeader
        onBack={() => history.push(`/stat/${userId}/all`)}
        title="祈愿历史"
        subTitle={`${userId}`}
      />
      <WishTable />
    </>
  );
};

const WishTable = () => {
  const [wishLogs, setWishLogs] = useState<WishLog[]>([]);
  const [gachaConfigs, setGachaConfigs] = useState<GachaConfig[]>([]);

  const { userId, gachaType } = useParams<{
    userId: string;
    gachaType: string;
  }>();

  useEffect(() => {
    const fetchLog = async () => {
      const gachaLog = await client.get<{ data: WishLog[] }>(`log/${userId}`, {
        params: {
          gachaType: gachaType === 'all' ? '' : gachaType,
        },
      });
      setWishLogs(gachaLog.data.data);
    };
    fetchLog();

    const fetchConfigs = async () => {
      const gachaLog = await client.get<{ data: GachaConfig[] }>('gacha');
      setGachaConfigs(gachaLog.data.data);
    };
    fetchConfigs();
  }, []);

  const columns = [
    {
      title: '卡池',
      filters: gachaConfigs.map(config => ({
        text: config.name,
        value: config.key,
      })),
      onFilter: (value: string | number | boolean, record: WishLog) =>
        record.gachaType === value,
      render: (value: WishLog) => {
        const config = gachaConfigs.find(
          config => config.key === value.gachaType,
        );
        if (config) {
          return config.name;
        }
      },
    },
    {
      title: '',
      dataIndex: ['Item'],
      render: (value: GachaItem) => (
        <Avatar itemInfo={value} rarity={value.rarity} />
      ),
    },
    { title: '物品名', dataIndex: ['Item', 'name'] },
    {
      title: '分类',
      dataIndex: ['Item', 'type'],
      filters: [
        { text: '角色', value: 'character' },
        { text: '武器', value: 'weapon' },
      ],
      onFilter: (value: string | number | boolean, record: WishLog) =>
        record.Item.type === value,
      render: (value: string) => (value === 'character' ? '角色' : '武器'),
    },
    {
      title: '稀有度',
      dataIndex: ['Item', 'rarity'],
      filters: [
        { text: '五星', value: '5' },
        { text: '四星', value: '4' },
        { text: '三星', value: '3' },
      ],
      onFilter: (value: string | number | boolean, record: WishLog) =>
        record.Item.rarity === value,
    },
    {
      title: '抽数',
      render: (value: WishLog) => {
        if (value.Item.rarity === '4') {
          return value.pityStar4;
        } else if (value.Item.rarity === '5') {
          return value.pityStar5;
        }
      },
    },
    { title: '日期', dataIndex: 'time' },
  ];

  return (
    <Table
      columns={columns}
      dataSource={wishLogs.map(wishLog => {
        if (wishLog.Item.type === 'character') {
          return {
            key: wishLog.id,
            ...wishLog,
            Item: {
              ...wishLog.Item,
              ...characters(wishLog.Item.id, queryOption as QueryOptions),
            },
          };
        } else {
          return {
            key: wishLog.id,
            ...wishLog,
            Item: {
              ...wishLog.Item,
              ...weapons(wishLog.Item.id, queryOption as QueryOptions),
            },
          };
        }
      })}
    />
  );
};
